package domain

import (
	"tibia-stats/repository"
)

func ListGuilds(guildRepository repository.GuildRepository) ([]string, error) {
	return guildRepository.ListGuilds()
}
