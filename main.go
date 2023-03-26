package main

import (
	"log"
	"tibia-exp-tracker/domain"
	"tibia-exp-tracker/dynamo"
	"tibia-exp-tracker/tibia"
)

func main() {
	expRepository, err := dynamo.InitializeExpRepository()
	if err != nil {
		log.Fatal(err)
	}
	guildMemberRepository, err := dynamo.InitializeGuildMembersRepository()
	if err != nil {
		log.Fatal(err)
	}

	apiClient := tibia.NewApiClient()

	err = domain.FetchExperience(apiClient, expRepository, "Peloria")
	if err != nil {
		log.Fatal(err)
	}

	err = domain.FetchGuildMembers(apiClient, guildMemberRepository, "Peloria")
	if err != nil {
		log.Fatal(err)
	}
}
