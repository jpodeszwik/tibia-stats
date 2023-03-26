package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-exp-tracker/domain"
	"tibia-exp-tracker/dynamo"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/slices"
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

	expHistory, err := domain.GetExperienceHistory(expRepository, event.PathParameters.PlayerName)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(expHistory, func(in repository.ExpHistory) ExpRecord {
		return ExpRecord{
			Exp:  in.Exp,
			Date: in.Date,
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
