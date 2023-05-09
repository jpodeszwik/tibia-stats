package domain

import "time"

type Guild struct {
	Name    string
	Members []GuildMember
	Date    string
}

type GuildMember struct {
	Name  string
	Level int
}

type GuildExp struct {
	GuildName    string
	HighScoreExp int64
	GainedExp    int64
	Date         time.Time
}

type GuildEvent struct {
	Name    string
	Members []GuildMember
}

type GuildMemberAction struct {
	GuildName     string
	Time          time.Time
	Level         int
	CharacterName string
	Action        Action
}
