package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/slices"
)

type LambdaEvent struct {
	PathParameters GetGuildHistoryEvent `json:"pathParameters"`
}

type GetGuildHistoryEvent struct {
	GuildName string `json:"guildName"`
}

type GuildHistoryRecord struct {
	PlayerName string `json:"playerName"`
	Date       string `json:"date"`
	Action     string `json:"action"`
	Level      int    `json:"level,omitempty"`
}

func HandleLambdaExecution(event LambdaEvent) ([]GuildHistoryRecord, error) {
	expRepository, err := dynamo.InitializeGuildMembersRepository()
	if err != nil {
		log.Fatal(err)
	}

	guildHistory, err := getGuildMemberHistory(expRepository, event.PathParameters.GuildName)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(guildHistory, func(in historyRecord) GuildHistoryRecord {
		return GuildHistoryRecord{
			Date:       in.Date,
			PlayerName: in.PlayerName,
			Action:     string(in.Action),
			Level:      in.Level,
		}
	}), nil
}

func main() {
	lambda.Start(HandleLambdaExecution)
}

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

type historyRecord struct {
	Date       string
	PlayerName string
	Action     Action
	Level      int
}

func getGuildMemberHistory(memberRepository *dynamo.GuildMemberRepository, guild string) ([]historyRecord, error) {
	guildHistory, err := memberRepository.GetGuildMembersHistory(guild, 7)
	if err != nil {
		return nil, err
	}

	ret := make([]historyRecord, 0)
	for i := 0; i < len(guildHistory)-1; i++ {
		records := getDiff(guildHistory[i], guildHistory[i+1])
		ret = append(ret, records...)
	}
	return ret, nil
}

func memberName(in domain.GuildMember) string {
	return in.Name
}

func getDiff(currentDay domain.Guild, previousDay domain.Guild) []historyRecord {
	currentMembers := NewStringSet(slices.MapSlice(currentDay.Members, memberName))
	previousMembers := NewStringSet(slices.MapSlice(previousDay.Members, memberName))

	ret := make([]historyRecord, 0)

	for _, member := range currentDay.Members {
		if !previousMembers.Contains(member.Name) {
			ret = append(ret, historyRecord{
				Date:       currentDay.Date,
				PlayerName: member.Name,
				Level:      member.Level,
				Action:     JOIN,
			})
		}
	}

	for _, member := range previousDay.Members {
		if !currentMembers.Contains(member.Name) {
			ret = append(ret, historyRecord{
				Date:       currentDay.Date,
				PlayerName: member.Name,
				Level:      member.Level,
				Action:     LEAVE,
			})
		}
	}

	return ret
}
