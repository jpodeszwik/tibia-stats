package domain

import (
	"log"
	"sync"
	"tibia-stats/repository"
	"tibia-stats/slices"
	"tibia-stats/tibia"
)

type fetchWorldResponse struct {
	world  string
	guilds []tibia.OverviewGuild
	err    error
}

func ETLGuildMembers(ac *tibia.ApiClient, guildRepository repository.GuildRepository, memberRepository repository.GuildMemberRepository) error {
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
	log.Printf("Fetching %v worlds with %v workers", len(worlds), workers)
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
	log.Printf("Finished world fetching")

	allGuilds := make([]tibia.OverviewGuild, 0)
	for worldResponse := range worldGuilds {
		if worldResponse.err != nil {
			log.Printf("Error fetchin world %v guilds", worldResponse.world)
			continue
		}
		allGuilds = append(allGuilds, worldResponse.guilds...)
	}

	err = guildRepository.StoreGuilds(slices.MapSlice(allGuilds, func(in tibia.OverviewGuild) string {
		return in.Name
	}))
	if err != nil {
		log.Printf("failed to store guilds %v", err)
	}

	log.Printf("Found %v guilds", len(allGuilds))
	allGuildsChan := make(chan string, 100)
	go func() {
		for _, guild := range allGuilds {
			allGuildsChan <- guild.Name
		}
		close(allGuildsChan)
	}()

	log.Printf("Fetching and storing %v guilds with %v workers", len(allGuilds), workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for guildName := range allGuildsChan {
				err := fetchAndStoreGuild(guildName, ac, memberRepository)
				if err != nil {
					log.Printf("Failed to store guild %s members %v", guildName, err)
				}
			}
		}()
	}
	wg.Wait()

	log.Printf("Done fetching and storing guilds")
	return nil
}

func fetchAndStoreGuild(guildName string, ac *tibia.ApiClient, mr repository.GuildMemberRepository) error {
	guild, err := ac.FetchGuild(guildName)
	if err != nil {
		return err
	}

	memberNames := slices.MapSlice(guild.Members, func(in tibia.GuildMemberResponse) repository.GuildMember {
		return repository.GuildMember{
			Name:  in.Name,
			Level: in.Level,
		}
	})

	return mr.StoreGuildMembers(guildName, memberNames)
}