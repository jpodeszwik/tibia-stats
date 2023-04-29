package domain

type Guild struct {
	Name    string
	Members []GuildMember
	Date    string
}

type GuildMember struct {
	Name  string
	Level int
}
