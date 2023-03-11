package main

import (
	_ "github.com/lib/pq"
	"log"
	"tibia-exp-tracker/actions"
	"tibia-exp-tracker/postgres"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/tibia"
)

func main() {
	db, err := postgres.InitializePostgresDb()
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.CloseDb(db)

	expRepository := repository.NewPostgresExpRepository(db)
	guildMemberRepository := repository.NewPostgresGuildMemberRepository(db)
	apiClient := tibia.NewApiClient()

	err = actions.FetchExperience(apiClient, expRepository, "Peloria")
	if err != nil {
		log.Fatal(err)
	}

	err = actions.FetchGuildMembers(apiClient, guildMemberRepository, "Peloria")
	if err != nil {
		log.Fatal(err)
	}
}
