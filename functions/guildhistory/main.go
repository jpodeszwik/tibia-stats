package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/slices"
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
	Level      int    `json:"level,omitempty"`
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
			Level:      in.Level,
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
