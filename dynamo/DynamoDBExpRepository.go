package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"tibia-stats/repository"
	"tibia-stats/slices"
)

const isotime = "2006-01-02"

type dynamoDBExpRepository struct {
	dynamoDB  *dynamodb.Client
	tableName string
}

func (d dynamoDBExpRepository) StoreExperiences(expData []repository.ExpData) error {
	chunks := calculateChunks(len(expData), 25)
	expDataChunks := slices.SplitSlice(expData, chunks)
	log.Printf("Chunks %v", len(expDataChunks))

	expDataChan := make(chan []repository.ExpData, chunks)
	for _, chunk := range expDataChunks {
		expDataChan <- chunk
	}
	close(expDataChan)

	ret := make(chan error, chunks)
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
	for i := 0; i < chunks; i++ {
		if i%100 == 0 && i != 0 {
			log.Printf("%v done", i)
		}
		err := <-ret
		if err != nil {
			log.Printf("Error storing %v", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v errors storing data", len(errs))
	}

	return nil
}

func (d dynamoDBExpRepository) GetExpHistory(name string, limit int) ([]repository.ExpHistory, error) {
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

	return slices.MapSlice(out.Items, func(in map[string]types.AttributeValue) repository.ExpHistory {
		m := make(map[string]string)
		err = attributevalue.UnmarshalMap(in, &m)

		return repository.ExpHistory{
			Date: m["date"],
			Exp:  m["exp"],
		}
	}), nil
}

func mapExpData(ed repository.ExpData) types.WriteRequest {
	m := map[string]interface{}{
		"playerName": ed.Name,
		"date":       ed.Date.Format(isotime),
		"exp":        ed.Exp,
	}

	marshalled, err := attributevalue.MarshalMap(m)
	if err != nil {
		log.Printf("Error marshalling json %v", err)
		return types.WriteRequest{}
	}

	return types.WriteRequest{
		PutRequest: &types.PutRequest{
			Item: marshalled,
		},
	}
}

func calculateChunks(count int, maxChunkSize int) int {
	chunks := count / maxChunkSize

	if count%maxChunkSize == 0 {
		return chunks
	}
	return chunks + 1
}

func NewDynamoDBExpRepository(db *dynamodb.Client, tableName string) repository.ExpRepository {
	return &dynamoDBExpRepository{dynamoDB: db, tableName: tableName}
}
