package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/formats"
	"tibia-stats/utils/slices"
)

type LambdaEvent struct {
	PathParameters GetExpEvent `json:"pathParameters"`
}

type GetExpEvent struct {
	GuildName string `json:"guildName"`
}

type ExpRecord struct {
	Date         string `json:"date"`
	HighScoreExp int64  `json:"exp"`
	GainedExp    int64  `json:"gainedExp"`
}

func HandleLambdaExecution(event LambdaEvent) ([]ExpRecord, error) {
	expRepository, err := dynamo.InitializeGuildExpRepository()
	if err != nil {
		log.Fatal(err)
	}

	expHistory, err := expRepository.GetExpHistory(event.PathParameters.GuildName, 30)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(expHistory, func(in domain.GuildExp) ExpRecord {
		return ExpRecord{
			HighScoreExp: in.HighScoreExp,
			GainedExp:    in.GainedExp,
			Date:         in.Date.Format(formats.IsoDate),
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
