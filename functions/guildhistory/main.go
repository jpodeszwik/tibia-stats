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
	guildMemberActionRepository, err := dynamo.InitializeGuildMemberActionRepository()
	if err != nil {
		log.Fatal(err)
	}

	actions, err := guildMemberActionRepository.GetActions(event.PathParameters.GuildName)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(actions, func(in domain.GuildMemberAction) GuildHistoryRecord {
		return GuildHistoryRecord{
			Date:       in.Time.Format(formats.IsoDate),
			PlayerName: in.CharacterName,
			Action:     string(in.Action),
			Level:      in.Level,
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
