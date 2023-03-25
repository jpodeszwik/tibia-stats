package dynamo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"os"
	"tibia-exp-tracker/repository"
)

func initializeDynamoDB() (client *dynamodb.Client, err error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func InitializeExpRepository() (repository.ExpRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		log.Printf("unable to initialize DynamoDB, %v", err)
		return nil, err
	}

	expTable, exists := os.LookupEnv("TIBIA_EXP_TABLE")
	if !exists {
		return nil, errors.New("TIBIA_EXP_TABLE not set")
	}

	return NewDynamoDBExpRepository(client, expTable), nil
}

func InitializeGuildMembersRepository() (repository.GuildMemberRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		log.Printf("unable to initialize DynamoDB, %v", err)
		return nil, err
	}

	guildMembersTable, exists := os.LookupEnv("TIBIA_GUILD_MEMBERS_TABLE")
	if !exists {
		return nil, errors.New("TIBIA_EXP_TABLE not set")
	}

	return NewDynamoDBGuildMemberRepository(client, guildMembersTable), nil
}
