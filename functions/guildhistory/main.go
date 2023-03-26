package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-exp-tracker/domain"
	"tibia-exp-tracker/dynamo"
	"tibia-exp-tracker/slices"
)

type LambdaEvent struct {
	PathParameters GetGuildHistoryEvent `json:"pathParameters"`
}

type GetGuildHistoryEvent struct {
	GuildName string `json:"guildName"`
}

type GuildHistoryRecord struct {
	PlayerName string `json:"playerName"`
	Date       string `json:"date"`
	Action     string `json:"action"`
}

func HandleLambdaExecution(event LambdaEvent) ([]GuildHistoryRecord, error) {
	expRepository, err := dynamo.InitializeGuildMembersRepository()
	if err != nil {
		log.Fatal(err)
	}

	guildHistory, err := domain.GetGuildMemberHistory(expRepository, event.PathParameters.GuildName)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(guildHistory, func(in domain.HistoryRecord) GuildHistoryRecord {
		return GuildHistoryRecord{
			Date:       in.Date,
			PlayerName: in.PlayerName,
			Action:     string(in.Action),
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}