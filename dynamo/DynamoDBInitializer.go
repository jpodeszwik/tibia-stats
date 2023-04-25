package dynamo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"os"
	"tibia-stats/repository"
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
		return nil, errors.New("TIBIA_GUILD_MEMBERS_TABLE not set")
	}

	return NewDynamoDBGuildMemberRepository(client, guildMembersTable), nil
}

func InitializeGuildRepository() (repository.GuildRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		log.Printf("unable to initialize DynamoDB, %v", err)
		return nil, err
	}

	guildsTable, exists := os.LookupEnv("TIBIA_GUILDS_TABLE")
	if !exists {
		return nil, errors.New("TIBIA_GUILDS_TABLE not set")
	}

	return NewDynamoDBGuildRepository(client, guildsTable), nil
}

func InitializeDeathRepository() (*DeathRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		return nil, err
	}

	deathtable, exists := os.LookupEnv("DEATH_TABLE_NAME")
	if !exists {
		return nil, errors.New("DEATH_TABLE_NAME not set")
	}

	characterNameIndex, exists := os.LookupEnv("DEATH_TABLE_CHARACTER_NAME_DATE_INDEX")
	if !exists {
		return nil, errors.New("DEATH_TABLE_CHARACTER_NAME_DATE_INDEX not set")
	}

	guildTimeIndex, exists := os.LookupEnv("DEATH_TABLE_GUILD_TIME_INDEX")
	if !exists {
		return nil, errors.New("DEATH_TABLE_GUILD_TIME_INDEX not set")
	}

	return NewDeathRepository(client, deathtable, characterNameIndex, guildTimeIndex), nil
}
