package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
	"tibia-stats/domain"
	"tibia-stats/utils/formats"
	"tibia-stats/utils/slices"
	"time"
)

type GuildMemberActionRepository struct {
	client             *dynamodb.Client
	tableName          string
	guildNameTimeIndex string
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

func (g *GuildMemberActionRepository) GetActions(guildName string) ([]domain.GuildMemberAction, error) {
	out, err := g.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(g.tableName),
		IndexName:        aws.String(g.guildNameTimeIndex),
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(30),
		KeyConditions: map[string]types.Condition{
			"guildName": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: guildName},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return []domain.GuildMemberAction{}, nil
	}

	return slices.MapSliceWithError(out.Items, func(in map[string]types.AttributeValue) (domain.GuildMemberAction, error) {
		m := make(map[string]string)
		err = attributevalue.UnmarshalMap(in, &m)
		if err != nil {
			return domain.GuildMemberAction{}, err
		}

		level, err := strconv.Atoi(m["level"])
		if err != nil {
			return domain.GuildMemberAction{}, err
		}

		parsedTime, err := time.Parse(formats.IsoDateTime, m["time"])
		if err != nil {
			return domain.GuildMemberAction{}, err
		}

		return domain.GuildMemberAction{
			GuildName:     m["guildName"],
			Time:          parsedTime,
			Level:         level,
			CharacterName: m["characterName"],
			Action:        domain.Action(m["action"]),
		}, nil
	})
}

func NewGuildMemberActionRepository(client *dynamodb.Client, tableName string, guildNameTimeIndex string) *GuildMemberActionRepository {
	return &GuildMemberActionRepository{
		client:             client,
		tableName:          tableName,
		guildNameTimeIndex: guildNameTimeIndex,
	}
}
