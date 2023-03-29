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
	apiClient := tibia.NewApiClient()

	err = domain.ETLGuildMembers(apiClient, guildMemberRepository)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
