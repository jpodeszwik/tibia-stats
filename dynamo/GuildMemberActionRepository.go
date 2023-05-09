package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"tibia-stats/domain"
	"tibia-stats/utils/formats"
)

type GuildMemberActionRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (g *GuildMemberActionRepository) StoreGuildMemberAction(ga domain.GuildMemberAction) error {
	formattedTime := ga.Time.Format(formats.IsoDateTime)
	m := map[string]interface{}{
		"guildName":          ga.GuildName,
		"time":               formattedTime,
		"action":             ga.Action,
		"characterName":      ga.CharacterName,
		"level":              ga.Level,
		"time-characterName": formattedTime + "-" + ga.CharacterName,
	}

	marshalled, err := attributevalue.MarshalMap(m)
	if err != nil {
		return err
	}

	_, err = g.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(g.tableName),
		Item:      marshalled,
	})
	return err
}

func NewGuildMemberActionRepository(client *dynamodb.Client, tableName string) *GuildMemberActionRepository {
	return &GuildMemberActionRepository{
		client:    client,
		tableName: tableName,
	}
}
