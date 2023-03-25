package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-exp-tracker/domain"
	"tibia-exp-tracker/dynamo"
	"tibia-exp-tracker/tibia"
)

func HandleLambdaExecution() {
	expRepository, err := dynamo.InitializeExpRepository()
	if err != nil {
		log.Fatal(err)
	}
	apiClient := tibia.NewApiClient()

	err = domain.FetchExperience(apiClient, expRepository, "Peloria")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
