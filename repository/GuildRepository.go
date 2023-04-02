package repository

type GuildRepository interface {
	StoreGuilds(guilds []string) error
	ListGuilds() ([]string, error)
}
