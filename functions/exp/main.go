package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
	"log"
	"tibia-exp-tracker/actions"
	"tibia-exp-tracker/postgres"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/tibia"
)

func HandleLambdaExecution() {
	db, err := postgres.InitializePostgresDb()
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.CloseDb(db)

	expRepository := repository.NewPostgresExpRepository(db)
	apiClient := tibia.NewApiClient()

	err = actions.FetchExperience(apiClient, expRepository, "Peloria")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
