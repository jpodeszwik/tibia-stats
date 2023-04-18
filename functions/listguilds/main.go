package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
)

func HandleLambdaExecution() ([]string, error) {
	guildRepository, err := dynamo.InitializeGuildRepository()
	if err != nil {
		log.Fatal(err)
	}

	return domain.ListGuilds(guildRepository)
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
