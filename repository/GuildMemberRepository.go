package repository

type Guild struct {
	Name    string
	Members []string
	Date    string
}

type GuildMemberRepository interface {
	StoreGuildMembers(guild string, members []string) error
	GetGuildsHistory(guild string, limit int) ([]Guild, error)
}
