package repository

type Guild struct {
	Name    string
	Members []GuildMember
	Date    string
}

type GuildMember struct {
	Name  string
	Level int
}

type GuildMemberRepository interface {
	StoreGuildMembers(guild string, members []GuildMember) error
	GetGuildMembersHistory(guild string, limit int) ([]Guild, error)
}
