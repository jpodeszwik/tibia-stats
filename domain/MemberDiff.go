package domain

import "tibia-stats/utils/slices"

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

type MemberDiffRecord struct {
	CharacterName string
	Action        Action
	Level         int
}

func memberName(in GuildMember) string {
	return in.Name
}

func MemberDiff(currentMembers []GuildMember, previousMembers []GuildMember) []MemberDiffRecord {
	currentMemberNames := NewStringSet(slices.MapSlice(currentMembers, memberName))
	previousMemberNames := NewStringSet(slices.MapSlice(previousMembers, memberName))

	ret := make([]MemberDiffRecord, 0)

	for _, member := range currentMembers {
		if !previousMemberNames.Contains(member.Name) {
			ret = append(ret, MemberDiffRecord{
				CharacterName: member.Name,
				Level:         member.Level,
				Action:        JOIN,
			})
		}
	}

	for _, member := range previousMembers {
		if !currentMemberNames.Contains(member.Name) {
			ret = append(ret, MemberDiffRecord{
				CharacterName: member.Name,
				Level:         member.Level,
				Action:        LEAVE,
			})
		}
	}

	return ret
}
