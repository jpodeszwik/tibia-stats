package tracker

import (
	"tibia-stats/dynamo"
	"tibia-stats/utils/logger"
)

type Guilds struct {
	gr *dynamo.GuildRepository
}

func (g *Guilds) Handle(guilds []string) {
	err := g.gr.StoreGuilds(guilds)
	if err != nil {
		logger.Error.Printf("%v", err)
	}
}

func NewGuilds(gr *dynamo.GuildRepository) *Guilds {
	return &Guilds{gr: gr}
}
