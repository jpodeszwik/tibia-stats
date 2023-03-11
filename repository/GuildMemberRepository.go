package repository

type GuildMemberRepository interface {
	StoreGuildMembers(guild string, members []string) error
}
