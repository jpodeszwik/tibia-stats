package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"tibia-stats/domain"
	"tibia-stats/utils/formats"
	"tibia-stats/utils/logger"
	"tibia-stats/utils/slices"
)

type ExpRepository struct {
	dynamoDB  *dynamodb.Client
	tableName string
}

func (d ExpRepository) StoreExperiences(expData []domain.ExpData) error {
	expDataChunks := slices.SplitSlice(expData, 25)
	logger.Info.Printf("Chunks %v", len(expDataChunks))

	expDataChan := make(chan []domain.ExpData, len(expDataChunks))
	for _, chunk := range expDataChunks {
		expDataChan <- chunk
	}
	close(expDataChan)

	ret := make(chan error, len(expDataChunks))
	defer close(ret)

	workers := 8
	for i := 0; i < workers; i++ {
		go func() {
			for chunk := range expDataChan {
				writeRequests := slices.MapSlice(chunk, mapExpData)
				_, err := d.dynamoDB.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						d.tableName: writeRequests,
					},
				})
				ret <- err
			}
		}()
	}

	errs := make([]error, 0)
	for i := 0; i < len(expDataChunks); i++ {
		if i%100 == 0 && i != 0 {
			logger.Info.Printf("%v done", i)
		}
		err := <-ret
		if err != nil {
			logger.Error.Printf("Error storing %v", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v errors storing data", len(errs))
	}

	return nil
}

func (d ExpRepository) GetExpHistory(name string, limit int) ([]domain.ExpHistory, error) {
	out, err := d.dynamoDB.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(d.tableName),
		IndexName:        aws.String("playerName-date-index"),
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
		KeyConditions: map[string]types.Condition{
			"playerName": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: name},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return slices.MapSlice(out.Items, func(in map[string]types.AttributeValue) domain.ExpHistory {
		m := make(map[string]string)
		err = attributevalue.UnmarshalMap(in, &m)

		return domain.ExpHistory{
			Date: m["date"],
			Exp:  m["exp"],
		}
	}), nil
}

func mapExpData(ed domain.ExpData) types.WriteRequest {
	m := map[string]interface{}{
		"playerName": ed.Name,
		"date":       ed.Date.Format(formats.IsoDate),
		"exp":        ed.Exp,
	}

	marshalled, err := attributevalue.MarshalMap(m)
	if err != nil {
		logger.Error.Printf("Error marshalling json %v", err)
		return types.WriteRequest{}
	}

	return types.WriteRequest{
		PutRequest: &types.PutRequest{
			Item: marshalled,
		},
	}
}

func NewExpRepository(db *dynamodb.Client, tableName string) *ExpRepository {
	return &ExpRepository{dynamoDB: db, tableName: tableName}
}
