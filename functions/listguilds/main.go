package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-stats/dynamo"
)

func HandleLambdaExecution() ([]string, error) {
	guildRepository, err := dynamo.InitializeGuildRepository()
	if err != nil {
		log.Fatal(err)
	}

	return guildRepository.ListGuilds()
}

func main() {
	lambda.Start(HandleLambdaExecution)
}
