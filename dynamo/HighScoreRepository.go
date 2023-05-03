package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"tibia-stats/domain"
	"tibia-stats/utils/formats"
	"time"
)

type HighScoreRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (hsr *HighScoreRepository) GetHighScore(worldName string, date time.Time) (*domain.WorldExperience, error) {
	out, err := hsr.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName: aws.String(hsr.tableName),
		Limit:     aws.Int32(1),
		KeyConditions: map[string]types.Condition{
			"worldName": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: worldName},
				},
			},
			"date": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: date.Format(formats.IsoDate)},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, nil
	}

	item := out.Items[0]
	worldValue := item["worldName"].(*types.AttributeValueMemberS).Value
	experienceValue := item["experience"].(*types.AttributeValueMemberM).Value

	experience := make(map[string]int64)
	err = attributevalue.UnmarshalMap(experienceValue, &experience)
	if err != nil {
		return nil, err
	}

	return &domain.WorldExperience{
		World:      worldValue,
		Experience: experience,
	}, nil
}

func (hsr *HighScoreRepository) StoreHighScore(highScore domain.WorldExperience) error {
	m := map[string]interface{}{
		"worldName":  highScore.World,
		"date":       time.Now().Format(formats.IsoDate),
		"experience": highScore.Experience,
	}

	marshalled, err := attributevalue.MarshalMap(m)
	if err != nil {
		return err
	}

	_, err = hsr.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(hsr.tableName),
		Item:      marshalled,
	})

	return err
}

func NewHighScoreRepository(client *dynamodb.Client, tableName string) *HighScoreRepository {
	return &HighScoreRepository{client: client, tableName: tableName}
}
