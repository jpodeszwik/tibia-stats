package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/slices"
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

	ret := make(chan error, chunks)
	defer close(ret)

	for _, chunk := range expDataChunks {
		go func(chunk2 []repository.ExpData) {
			writeRequests := slices.MapSlice(chunk2, mapExpData)
			_, err := d.dynamoDB.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{
					d.tableName: writeRequests,
				},
			})
			ret <- err
		}(chunk)
	}

	var err error
	for i := 0; i < chunks; i++ {
		err = <-ret
		log.Printf("%v done %v ", i+1, err)
	}

	return err
}

func (d dynamoDBExpRepository) GetExpHistory(name string, limit int) ([]repository.ExpHistory, error) {
	out, err := d.dynamoDB.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(d.tableName),
		IndexName:        aws.String("playerName-date-index"),
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
		KeyConditions: map[string]types.Condition{
			"playerName": types.Condition{
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
