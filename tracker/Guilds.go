package tracker

import (
	"log"
	"tibia-stats/dynamo"
)

type Guilds struct {
	gr *dynamo.GuildRepository
}

func (g *Guilds) Handle(guilds []string) {
	err := g.gr.StoreGuilds(guilds)
	if err != nil {
		log.Printf("%v", err)
	}
}

func NewGuilds(gr *dynamo.GuildRepository) *Guilds {
	return &Guilds{gr: gr}
}
