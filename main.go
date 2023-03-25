package main

import (
	_ "github.com/lib/pq"
	"log"
	"tibia-exp-tracker/domain"
	"tibia-exp-tracker/postgres"
	"tibia-exp-tracker/tibia"
)

func main() {
	db, err := postgres.InitializePostgresDb()
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.CloseDb(db)

	expRepository := postgres.NewPostgresExpRepository(db)
	guildMemberRepository := postgres.NewPostgresGuildMemberRepository(db)
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
