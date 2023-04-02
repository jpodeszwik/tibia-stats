package domain

import (
	"tibia-exp-tracker/repository"
)

func ListGuilds(guildRepository repository.GuildRepository) ([]string, error) {
	return guildRepository.ListGuilds()
}
