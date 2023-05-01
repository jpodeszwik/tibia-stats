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

type GuildExpRepository struct {
	client             *dynamodb.Client
	tableName          string
	guildNameDateIndex string
}

func (ger *GuildExpRepository) StoreGuildExp(ge domain.GuildExp) error {
	m := map[string]interface{}{
		"guildName": ge.GuildName,
		"date":      ge.Date.Format(formats.IsoDate),
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

func (ger *GuildExpRepository) GetExpHistory(guildName string, limit int) ([]domain.GuildExp, error) {
	out, err := ger.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(ger.tableName),
		IndexName:        aws.String(ger.guildNameDateIndex),
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
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

	return slices.MapSlice(out.Items, func(in map[string]types.AttributeValue) domain.GuildExp {
		m := make(map[string]string)
		err = attributevalue.UnmarshalMap(in, &m)
		parsedDate, _ := time.Parse(formats.IsoDate, m["date"])
		exp, _ := strconv.ParseInt(m["exp"], 10, 64)

		return domain.GuildExp{
			Date:      parsedDate,
			Exp:       exp,
			GuildName: m["guildName"],
		}
	}), nil
}

func NewGuildExpRepository(client *dynamodb.Client, tableName string, guildNameDateIndex string) *GuildExpRepository {
	return &GuildExpRepository{
		client:             client,
		tableName:          tableName,
		guildNameDateIndex: guildNameDateIndex,
	}
}
