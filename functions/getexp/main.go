package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/slices"
)

type LambdaEvent struct {
	PathParameters GetExpEvent `json:"pathParameters"`
}

type GetExpEvent struct {
	PlayerName string `json:"playerName"`
}

type ExpRecord struct {
	Date string `json:"date"`
	Exp  string `json:"exp"`
}

func HandleLambdaExecution(event LambdaEvent) ([]ExpRecord, error) {
	expRepository, err := dynamo.InitializeExpRepository()
	if err != nil {
		log.Fatal(err)
	}

	expHistory, err := expRepository.GetExpHistory(event.PathParameters.PlayerName, 30)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(expHistory, func(in domain.ExpHistory) ExpRecord {
		return ExpRecord{
			Exp:  in.Exp,
			Date: in.Date,
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
