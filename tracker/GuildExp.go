package tracker

import (
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/formats"
)

type GuildExp struct {
	lastExp    map[string]*domain.GuildExp
	repository *dynamo.GuildExpRepository
}

func (ge *GuildExp) Handle(exp domain.GuildExp) {
	last, exists := ge.lastExp[exp.GuildName]

	if !exists || last.Exp != exp.Exp || last.Date.Format(formats.IsoDate) != exp.Date.Format(formats.IsoDate) {
		log.Printf("Storing %v", exp)
		ge.repository.StoreGuildExp(exp)
		ge.lastExp[exp.GuildName] = &exp
	}
}

func NewGuildExp(repository *dynamo.GuildExpRepository) *GuildExp {
	return &GuildExp{
		lastExp:    make(map[string]*domain.GuildExp),
		repository: repository,
	}
}
