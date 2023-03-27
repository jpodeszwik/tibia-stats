package domain

import (
	"log"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/slices"
	"tibia-exp-tracker/tibia"
)

func FetchGuildMembers(ac *tibia.ApiClient, memberRepository repository.GuildMemberRepository, world string) error {
	log.Printf("Fetching guilds for %v", world)
	guilds, err := ac.FetchGuilds(world)
	if err != nil {
		return err
	}

	log.Printf("Found %d guilds", len(guilds))
	chunks := 8
	guildChunks := slices.SplitSlice(guilds, chunks)
	res := make(chan []error, chunks)
	defer close(res)

	for _, slice := range guildChunks {
		go func(ogs []tibia.OverviewGuild) {
			res <- fetchAndStoreGuilds(ac, memberRepository, ogs)
		}(slice)
	}

	for i := 0; i < chunks; i++ {
		errors := <-res
		if len(errors) > 0 {
			log.Printf("Errors when processing guilds %v", errors)
		}
	}

	log.Printf("Done")
	return nil
}

func fetchAndStoreGuilds(ac *tibia.ApiClient, mr repository.GuildMemberRepository, guilds []tibia.OverviewGuild) []error {
	var errors = make([]error, 0)
	for _, overviewGuild := range guilds {
		err := fetchAndStoreGuild(overviewGuild.Name, ac, mr)
		if err != nil {
			log.Printf("Failed to store guild %s members %v", overviewGuild.Name, err)
			errors = append(errors, err)
		}
	}
	return errors
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
