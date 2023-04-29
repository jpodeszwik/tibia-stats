package dynamo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strings"
	"tibia-stats/domain"
	"tibia-stats/slices"
	"time"
)

type GuildMemberRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (d *GuildMemberRepository) GetGuildMembersHistory(guild string, limit int) ([]domain.Guild, error) {
	return d.queryGuild(&dynamodb.QueryInput{
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
}

func (d *GuildMemberRepository) queryGuild(queryInput *dynamodb.QueryInput) ([]domain.Guild, error) {
	out, err := d.client.Query(context.Background(), queryInput)

	if err != nil {
		return nil, err
	}

	return slices.MapSliceWithError(out.Items, func(in map[string]types.AttributeValue) (domain.Guild, error) {
		m := make(map[string]interface{})
		err = attributevalue.UnmarshalMap(in, &m)

		guildName, ok := m["guildName"].(string)
		if !ok {
			return domain.Guild{}, errors.New("failed to deserialize guild name")
		}

		date, ok := m["date"].(string)
		if !ok {
			return domain.Guild{}, errors.New("failed to deserialize date")
		}

		membersInt, ok := m["members"].([]interface{})
		if !ok {
			return domain.Guild{}, errors.New("failed to deserialize members")
		}
		members, err := slices.MapSliceWithError(membersInt, func(in interface{}) (domain.GuildMember, error) {
			s, ok := in.(string)
			if ok {
				return domain.GuildMember{
					Name: s,
				}, nil
			}

			m, ok := in.(map[string]interface{})
			if !ok {
				return domain.GuildMember{}, errors.New("failed to deserialize member")
			}
			name, ok := m["name"].(string)
			if !ok {
				return domain.GuildMember{}, errors.New("failed to deserialize member")
			}

			level, ok := m["level"].(float64)
			if !ok {
				return domain.GuildMember{
					Name: name,
				}, nil
			}

			return domain.GuildMember{
				Name:  name,
				Level: int(level),
			}, nil
		})

		if err != nil {
			return domain.Guild{}, err
		}

		return domain.Guild{
			Name:    guildName,
			Members: members,
			Date:    date,
		}, nil
	})
}

func (d *GuildMemberRepository) StoreGuildMembers(guild string, members []domain.GuildMember) error {
	mem := slices.MapSlice(members, func(in domain.GuildMember) map[string]interface{} {
		return map[string]interface{}{
			"name":  in.Name,
			"level": in.Level,
		}
	})

	m := map[string]interface{}{
		"guildName":      guild,
		"lowerGuildName": strings.ToLower(guild),
		"date":           time.Now().Format(isotime),
		"members":        mem,
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

func NewGuildMemberRepository(client *dynamodb.Client, tableName string) *GuildMemberRepository {
	return &GuildMemberRepository{client: client, tableName: tableName}
}
