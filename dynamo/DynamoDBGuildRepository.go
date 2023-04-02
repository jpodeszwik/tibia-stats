package dynamo

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/slices"
	"time"
)

type dynamoDBGuildRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (d *dynamoDBGuildRepository) ListGuilds() ([]string, error) {
	out, err := d.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName: aws.String(d.tableName),
		KeyConditions: map[string]types.Condition{
			"date": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: time.Now().Add(-24 * time.Hour).Format(isotime)},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	items := out.Items
	for _, item := range items {
		m := make(map[string]interface{})
		err = attributevalue.UnmarshalMap(item, &m)
		if err != nil {
			return nil, err
		}
		log.Printf("%v", reflect.TypeOf(m["guilds"]))
		guilds, ok := m["guilds"].([]interface{})
		if !ok {
			return nil, errors.New("failed to deserialize guilds")
		}
		return slices.MapSlice(guilds, func(in interface{}) string {
			return fmt.Sprintf("%v", in)
		}), nil
	}
	return nil, errors.New("guild not found")
}

func (d *dynamoDBGuildRepository) StoreGuilds(guilds []string) error {
	m := map[string]interface{}{
		"guilds": guilds,
		"date":   time.Now().Format(isotime),
	}

	marshalled, err := attributevalue.MarshalMap(m)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      marshalled,
	})

	return err
}

func NewDynamoDBGuildRepository(client *dynamodb.Client, tableName string) repository.GuildRepository {
	return &dynamoDBGuildRepository{client: client, tableName: tableName}
}
