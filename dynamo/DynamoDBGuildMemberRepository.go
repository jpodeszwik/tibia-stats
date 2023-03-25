package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"tibia-exp-tracker/repository"
	"time"
)

type dynamoDBGuildMemberRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (d *dynamoDBGuildMemberRepository) StoreGuildMembers(guild string, members []string) error {
	m := map[string]interface{}{
		"guildName": guild,
		"date":      time.Now().Format(isotime),
		"members":   members,
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

func NewDynamoDBGuildMemberRepository(client *dynamodb.Client, tableName string) repository.GuildMemberRepository {
	return &dynamoDBGuildMemberRepository{client: client, tableName: tableName}
}
