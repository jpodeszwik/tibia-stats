package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"tibia-stats/domain"
)

type GuildExpRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (ger *GuildExpRepository) StoreGuildExp(ge domain.GuildExp) error {
	m := map[string]interface{}{
		"guildName": ge.GuildName,
		"date":      ge.Date.Format(isotime),
		"exp":       ge.Exp,
	}

	marshalled, err := attributevalue.MarshalMap(m)
	if err != nil {
		return err
	}

	_, err = ger.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(ger.tableName),
		Item:      marshalled,
	})
	return err
}

func NewGuildExpRepository(client *dynamodb.Client, tableName string) *GuildExpRepository {
	return &GuildExpRepository{
		client:    client,
		tableName: tableName,
	}
}
