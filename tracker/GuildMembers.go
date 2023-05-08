package tracker

import (
	"tibia-stats/domain"
	"tibia-stats/dynamo"
)

type GuildMembers struct {
	guildMemberRepository *dynamo.GuildMemberRepository
}

func (gm *GuildMembers) HandleGuild(guild domain.GuildEvent) {
	gm.guildMemberRepository.StoreGuildMembers(guild.Name, guild.Members)
}

func NewGuildMembers(guildMemberRepository *dynamo.GuildMemberRepository) *GuildMembers {
	return &GuildMembers{
		guildMemberRepository: guildMemberRepository,
	}
}
