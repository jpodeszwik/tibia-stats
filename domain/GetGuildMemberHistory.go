package domain

import (
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/slices"
)

type Action string

const (
	JOIN  Action = "JOIN"
	LEAVE Action = "LEAVE"
)

type StringSet struct {
	data map[string]bool
}

func (s StringSet) Contains(value string) bool {
	_, ok := s.data[value]
	return ok
}

func NewStringSet(values []string) StringSet {
	data := make(map[string]bool)
	for _, value := range values {
		data[value] = true
	}
	return StringSet{
		data: data,
	}
}

type HistoryRecord struct {
	Date       string
	PlayerName string
	Action     Action
	Level      int
}

func GetGuildMemberHistory(memberRepository repository.GuildMemberRepository, guild string) ([]HistoryRecord, error) {
	guildHistory, err := memberRepository.GetGuildMembersHistory(guild, 30)
	if err != nil {
		return nil, err
	}

	ret := make([]HistoryRecord, 0)
	for i := 0; i < len(guildHistory)-1; i++ {
		records := getDiff(guildHistory[i], guildHistory[i+1])
		ret = append(ret, records...)
	}
	return ret, nil
}

func memberName(in repository.GuildMember) string {
	return in.Name
}

func getDiff(currentDay repository.Guild, previousDay repository.Guild) []HistoryRecord {
	currentMembers := NewStringSet(slices.MapSlice(currentDay.Members, memberName))
	previousMembers := NewStringSet(slices.MapSlice(previousDay.Members, memberName))

	ret := make([]HistoryRecord, 0)

	for _, member := range currentDay.Members {
		if !previousMembers.Contains(member.Name) {
			ret = append(ret, HistoryRecord{
				Date:       currentDay.Date,
				PlayerName: member.Name,
				Level:      member.Level,
				Action:     JOIN,
			})
		}
	}

	for _, member := range previousDay.Members {
		if !currentMembers.Contains(member.Name) {
			ret = append(ret, HistoryRecord{
				Date:       currentDay.Date,
				PlayerName: member.Name,
				Level:      member.Level,
				Action:     LEAVE,
			})
		}
	}

	return ret
}
