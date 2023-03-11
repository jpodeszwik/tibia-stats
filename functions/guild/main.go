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

	guildMemberRepository := repository.NewPostgresGuildMemberRepository(db)
	apiClient := tibia.NewApiClient()

	err = actions.FetchGuildMembers(apiClient, guildMemberRepository, "Peloria")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
