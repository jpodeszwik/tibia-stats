package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-exp-tracker/domain"
	"tibia-exp-tracker/dynamo"
	"tibia-exp-tracker/tibia"
)

func HandleLambdaExecution() {
	guildMemberRepository, err := dynamo.InitializeGuildMembersRepository()
	if err != nil {
		log.Fatal(err)
	}
	guildRepository, err := dynamo.InitializeGuildRepository()
	if err != nil {
		log.Fatal(err)
	}
	apiClient := tibia.NewApiClient()

	err = domain.ETLGuildMembers(apiClient, guildRepository, guildMemberRepository)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
