package tracker

import (
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/logger"
	"time"
)

type GuildMemberAction struct {
	guildMemberRepository       *dynamo.GuildMemberRepository
	guildMemberActionRepository *dynamo.GuildMemberActionRepository
}

func (gm *GuildMemberAction) HandleGuild(guild domain.GuildEvent) {
	date := time.Now()
	for {
		err := gm.handleGuild(guild, date)
		if err == nil {
			return
		} else {
			logger.Error.Printf("Failed to handle guild %v", err)
		}
	}
}

func (gm *GuildMemberAction) handleGuild(guild domain.GuildEvent, date time.Time) error {
	lastMembers, err := gm.guildMemberRepository.GetLastGuildMembers(guild.Name)
	if err != nil {
		return err
	}

	if lastMembers != nil {
		diff := domain.MemberDiff(guild.Members, lastMembers)
		for _, record := range diff {
			err := gm.guildMemberActionRepository.StoreGuildMemberAction(domain.GuildMemberAction{
				GuildName:     guild.Name,
				Time:          date,
				Level:         record.Level,
				CharacterName: record.CharacterName,
				Action:        record.Action,
			})
			if err != nil {
				return err
			}
		}
	}

	return gm.guildMemberRepository.StoreLastGuildMembers(guild.Name, guild.Members)
}

func NewGuildMemberAction(guildMemberRepository *dynamo.GuildMemberRepository, guildMemberActionRepository *dynamo.GuildMemberActionRepository) *GuildMemberAction {
	return &GuildMemberAction{
		guildMemberRepository:       guildMemberRepository,
		guildMemberActionRepository: guildMemberActionRepository,
	}
}
