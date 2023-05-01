package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/slices"
)

type LambdaEvent struct {
	PathParameters GetExpEvent `json:"pathParameters"`
}

type GetExpEvent struct {
	GuildName string `json:"guildName"`
}

type Death struct {
	CharacterName string `json:"characterName"`
	Time          string `json:"time"`
	Reason        string `json:"reason"`
}

func HandleLambdaExecution(event LambdaEvent) ([]Death, error) {
	deathRepository, err := dynamo.InitializeDeathRepository()
	if err != nil {
		log.Fatal(err)
	}

	deaths, err := deathRepository.GetGuildDeaths(event.PathParameters.GuildName)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(deaths, func(in domain.Death) Death {
		return Death{
			CharacterName: in.CharacterName,
			Time:          in.Time.Format("2006-01-02T15:04:05Z"),
			Reason:        in.Reason,
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
