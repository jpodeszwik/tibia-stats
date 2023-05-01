package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"sync"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"tibia-stats/utils/slices"
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

	err = etlGuildMembers(apiClient, guildRepository, guildMemberRepository)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleLambdaExecution)
}

type fetchWorldResponse struct {
	world  string
	guilds []tibia.OverviewGuild
	err    error
}

func etlGuildMembers(ac *tibia.ApiClient, guildRepository *dynamo.GuildRepository, memberRepository *dynamo.GuildMemberRepository) error {
	worlds, err := ac.FetchWorlds()
	if err != nil {
		return err
	}

	worldJobs := make(chan string, len(worlds))
	worldGuilds := make(chan fetchWorldResponse, len(worlds))

	for _, world := range worlds {
		worldJobs <- world.Name
	}
	close(worldJobs)

	workers := 8
	logger.Info.Printf("Fetching %v worlds with %v workers", len(worlds), workers)
	wg := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range worldJobs {
				guilds, err := ac.FetchGuilds(job)
				worldGuilds <- fetchWorldResponse{
					guilds: guilds,
					err:    err,
				}
			}
		}()
	}
	wg.Wait()
	close(worldGuilds)
	logger.Info.Printf("Finished world fetching")

	allGuilds := make([]tibia.OverviewGuild, 0)
	for worldResponse := range worldGuilds {
		if worldResponse.err != nil {
			logger.Error.Printf("Error fetching world %v guilds", worldResponse.world)
			continue
		}
		allGuilds = append(allGuilds, worldResponse.guilds...)
	}

	logger.Info.Printf("Found %v guilds", len(allGuilds))
	allGuildsChan := make(chan string, 100)
	go func() {
		for _, guild := range allGuilds {
			allGuildsChan <- guild.Name
		}
		close(allGuildsChan)
	}()

	logger.Info.Printf("Fetching and storing %v guilds with %v workers", len(allGuilds), workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for guildName := range allGuildsChan {
				err := fetchAndStoreGuild(guildName, ac, memberRepository)
				if err != nil {
					logger.Error.Printf("Failed to store guild %s members %v", guildName, err)
				}
			}
		}()
	}
	wg.Wait()

	logger.Info.Printf("Done fetching and storing guilds")
	return nil
}

func fetchAndStoreGuild(guildName string, ac *tibia.ApiClient, mr *dynamo.GuildMemberRepository) error {
	guild, err := ac.FetchGuild(guildName)
	if err != nil {
		return err
	}

	memberNames := slices.MapSlice(guild.Members, func(in tibia.GuildMemberResponse) domain.GuildMember {
		return domain.GuildMember{
			Name:  in.Name,
			Level: in.Level,
		}
	})

	return mr.StoreGuildMembers(guildName, memberNames)
}
