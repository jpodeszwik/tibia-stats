package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"tibia-stats/domain"
	"tibia-stats/utils/formats"
	"tibia-stats/utils/logger"
	"tibia-stats/utils/slices"
	"time"
)

type DeathRepository struct {
	client                 *dynamodb.Client
	tableName              string
	characterNameDateIndex string
	guildTimeIndex         string
}

func NewDeathRepository(client *dynamodb.Client, tableName string, characterNameDateIndex string, guildTimeIndex string) *DeathRepository {
	return &DeathRepository{
		client:                 client,
		tableName:              tableName,
		characterNameDateIndex: characterNameDateIndex,
		guildTimeIndex:         guildTimeIndex,
	}
}

func (dr *DeathRepository) StoreDeaths(deaths []domain.Death) error {
	var mapped []map[string]types.AttributeValue
	for _, death := range deaths {
		m := map[string]interface{}{
			"characterName": death.CharacterName,
			"time":          death.Time.Format(formats.IsoDateTime),
			"reason":        death.Reason,
		}
		if death.Guild != "" {
			m["guild"] = death.Guild
		}

		marshalled, err := attributevalue.MarshalMap(m)
		if err != nil {
			return err
		}
		mapped = append(mapped, marshalled)
	}

	chunks := slices.SplitSlice(mapped, 25)
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		} else if len(chunk) == 1 {
			_, err := dr.client.PutItem(context.Background(), &dynamodb.PutItemInput{
				TableName: aws.String(dr.tableName),
				Item:      mapped[0],
			})
			if err != nil {
				return err
			}
		} else {
			writeRequests := slices.MapSlice(chunk, toWriteRequest)
			_, err := dr.client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{
					dr.tableName: writeRequests,
				},
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func toWriteRequest(item map[string]types.AttributeValue) types.WriteRequest {
	return types.WriteRequest{
		PutRequest: &types.PutRequest{
			Item: item,
		},
	}
}

func (dr *DeathRepository) GetLastDeath(characterName string) (*domain.Death, error) {
	out, err := dr.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(dr.tableName),
		IndexName:        aws.String(dr.characterNameDateIndex),
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(1)),
		KeyConditions: map[string]types.Condition{
			"characterName": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: characterName},
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

	m := make(map[string]string)
	err = attributevalue.UnmarshalMap(out.Items[0], &m)

	parsedTime, err := time.Parse("2006-01-02T15:04:05Z", m["time"])
	if err != nil {
		logger.Error.Printf("Failed to parse time %v", err)
		return nil, err
	}

	return &domain.Death{
		CharacterName: m["characterName"],
		Guild:         m["guild"],
		Time:          parsedTime,
		Reason:        m["reason"],
	}, nil
}

func (dr *DeathRepository) GetGuildDeaths(guildName string) ([]domain.Death, error) {
	out, err := dr.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(dr.tableName),
		IndexName:        aws.String(dr.guildTimeIndex),
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(30)),
		KeyConditions: map[string]types.Condition{
			"guild": {
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

	var deaths []domain.Death
	for _, item := range out.Items {
		m := make(map[string]string)
		err = attributevalue.UnmarshalMap(item, &m)

		parsedTime, err := time.Parse(formats.IsoDateTime, m["time"])
		if err != nil {
			logger.Error.Printf("Failed to parse time %v", err)
			return nil, err
		}

		deaths = append(deaths, domain.Death{
			CharacterName: m["characterName"],
			Guild:         m["guild"],
			Time:          parsedTime,
			Reason:        m["reason"],
		})
	}

	return deaths, nil
}
