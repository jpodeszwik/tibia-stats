package dynamo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
	"tibia-stats/utils/logger"
)

func initializeDynamoDB() (client *dynamodb.Client, err error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func InitializeExpRepository() (*ExpRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		logger.Error.Printf("unable to initialize DynamoDB, %v", err)
		return nil, err
	}

	expTable, exists := os.LookupEnv("TIBIA_EXP_TABLE")
	if !exists {
		return nil, errors.New("TIBIA_EXP_TABLE not set")
	}

	return NewExpRepository(client, expTable), nil
}

func InitializeGuildMembersRepository() (*GuildMemberRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		logger.Error.Printf("unable to initialize DynamoDB, %v", err)
		return nil, err
	}

	guildMembersTable, exists := os.LookupEnv("TIBIA_GUILD_MEMBERS_TABLE")
	if !exists {
		return nil, errors.New("TIBIA_GUILD_MEMBERS_TABLE not set")
	}

	return NewGuildMemberRepository(client, guildMembersTable), nil
}

func InitializeGuildRepository() (*GuildRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		logger.Error.Printf("unable to initialize DynamoDB, %v", err)
		return nil, err
	}

	guildsTable, exists := os.LookupEnv("TIBIA_GUILDS_TABLE")
	if !exists {
		return nil, errors.New("TIBIA_GUILDS_TABLE not set")
	}

	return NewGuildRepository(client, guildsTable), nil
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

func InitializeGuildExpRepository() (*GuildExpRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		return nil, err
	}

	guildExpTable, exists := os.LookupEnv("GUILD_EXP_TABLE_NAME")
	if !exists {
		return nil, errors.New("GUILD_EXP_TABLE_NAME not set")
	}

	guildNameDateIndex, exists := os.LookupEnv("GUILD_EXP_GUILD_NAME_DATE_INDEX")
	if !exists {
		return nil, errors.New("GUILD_EXP_GUILD_NAME_DATE_INDEX not set")
	}

	return NewGuildExpRepository(client, guildExpTable, guildNameDateIndex), nil
}

func InitializeHighScoreRepository() (*HighScoreRepository, error) {
	client, err := initializeDynamoDB()
	if err != nil {
		return nil, err
	}

	highScoreTable, exists := os.LookupEnv("HIGHSCORE_TABLE")
	if !exists {
		return nil, errors.New("HIGHSCORE_TABLE not set")
	}

	return NewHighScoreRepository(client, highScoreTable), nil
}
