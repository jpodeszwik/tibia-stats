package dynamo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/slices"
	"time"
)

type dynamoDBGuildMemberRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (d *dynamoDBGuildMemberRepository) GetGuildsHistory(guild string, limit int) ([]repository.Guild, error) {
	out, err := d.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(d.tableName),
		IndexName:        aws.String("guildName-date-index"),
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
		KeyConditions: map[string]types.Condition{
			"guildName": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: guild},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return slices.MapSliceWithError(out.Items, func(in map[string]types.AttributeValue) (repository.Guild, error) {
		m := make(map[string]interface{})
		err = attributevalue.UnmarshalMap(in, &m)

		guildName, ok := m["guildName"].(string)
		if !ok {
			return repository.Guild{}, errors.New("failed to deserialize guild name")
		}

		date, ok := m["date"].(string)
		if !ok {
			return repository.Guild{}, errors.New("failed to deserialize date")
		}

		membersInt, ok := m["members"].([]interface{})
		if !ok {
			return repository.Guild{}, errors.New("failed to deserialize members")
		}
		members, err := slices.MapSliceWithError(membersInt, func(in interface{}) (string, error) {
			s, ok := in.(string)
			if !ok {
				return s, errors.New("failed to deserialize member")
			}
			return s, nil
		})

		if err != nil {
			return repository.Guild{}, err
		}

		return repository.Guild{
			Name:    guildName,
			Members: members,
			Date:    date,
		}, nil
	})
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