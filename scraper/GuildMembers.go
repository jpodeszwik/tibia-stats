package scraper

import (
	"sync"
	"tibia-stats/domain"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"tibia-stats/utils/slices"
	"time"
)

const guildMembersRefreshInterval = 2 * time.Hour

type GuildMembers struct {
	api               *tibia.ApiClient
	guilds            *Guilds
	playerGuild       map[string]string
	handler           Handler[map[string]string]
	guildEventHandler Handler[domain.GuildEvent]
	m                 sync.RWMutex
}

type workResult struct {
	guildName string
	members   []domain.GuildMember
	err       error
}

func (gm *GuildMembers) Start() {
	ticker := time.NewTicker(guildMembersRefreshInterval)
	err := gm.initialFetch(4)
	if err != nil {
		logger.Error.Fatal(err)
	}

	go func() {
		for range ticker.C {
			gm.fetchGuildMembers()
		}
	}()
}

func (gm *GuildMembers) initialFetch(workers int) error {
	start := time.Now()

	var wg sync.WaitGroup

	work := make(chan string, 100)
	result := make(chan workResult, 100)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for guildName := range work {
				guild, err := retry(func() (*tibia.GuildResponse, error) {
					return gm.api.FetchGuild(guildName)
				}, 5)
				if err != nil {
					logger.Error.Printf("Failed to fetch guild %v %v", guildName, err)
					result <- workResult{
						guildName: guildName,
						err:       err,
					}
				}
				var members []domain.GuildMember
				for _, member := range guild.Members {
					members = append(members, domain.GuildMember{
						Name:  member.Name,
						Level: member.Level,
					})
				}
				result <- workResult{
					guildName: guildName,
					members:   members,
				}
			}
		}()
	}

	guilds := gm.guilds.getGuilds()
	go func() {
		for _, guildName := range guilds {
			work <- guildName
		}
		close(work)
		wg.Wait()
		close(result)
	}()

	playerGuild := make(map[string]string)
	var err error
	for res := range result {
		if res.err != nil {
			err = res.err
			// even though this will fail whole job we need to read all results to let workers finish
			continue
		}

		gm.guildEventHandler(domain.GuildEvent{
			Name:    res.guildName,
			Members: res.members,
		})

		for _, member := range res.members {
			playerGuild[member.Name] = res.guildName
		}
	}
	if err != nil {
		return err
	}

	gm.m.Lock()
	gm.playerGuild = playerGuild
	gm.m.Unlock()
	gm.handler(playerGuild)

	logger.Info.Printf("Finished fetching %v guilds members with %v workers, %v memberships found in %v", len(guilds), workers, len(playerGuild), time.Since(start))
	return nil
}

func (gm *GuildMembers) fetchGuildMembers() error {
	start := time.Now()
	guilds := gm.guilds.getGuilds()

	playerGuild := make(map[string]string)
	for _, guildName := range guilds {
		guild, err := retry(func() (*tibia.GuildResponse, error) {
			return gm.api.FetchGuild(guildName)
		}, 5)
		if err != nil {
			return err
		}

		gm.guildEventHandler(domain.GuildEvent{
			Name: guild.Name,
			Members: slices.MapSlice(guild.Members, func(in tibia.GuildMemberResponse) domain.GuildMember {
				return domain.GuildMember{
					Name:  in.Name,
					Level: in.Level,
				}
			}),
		})

		for _, member := range guild.Members {
			playerGuild[member.Name] = guildName
		}
	}

	gm.m.Lock()
	gm.playerGuild = playerGuild
	gm.m.Unlock()
	gm.handler(playerGuild)

	logger.Info.Printf("Finished fetching %v guilds for members, %v memberships found in %v", len(guilds), len(playerGuild), time.Since(start))
	return nil
}

func (gm *GuildMembers) GetPlayerGuild() map[string]string {
	gm.m.RLock()
	defer gm.m.RUnlock()
	return gm.playerGuild
}

func NewGuildMembers(client *tibia.ApiClient, guilds *Guilds, handler Handler[map[string]string], guildEventHandler Handler[domain.GuildEvent]) *GuildMembers {
	return &GuildMembers{
		api:               client,
		guilds:            guilds,
		playerGuild:       make(map[string]string),
		handler:           handler,
		guildEventHandler: guildEventHandler,
	}
}
