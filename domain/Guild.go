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
