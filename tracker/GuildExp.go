package tracker

import (
	"log"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
)

const isoDate = "2006-01-02"

type GuildExp struct {
	lastExp    map[string]*domain.GuildExp
	repository *dynamo.GuildExpRepository
}

func (ge *GuildExp) Handle(exp domain.GuildExp) {
	last, exists := ge.lastExp[exp.GuildName]

	if !exists || last.Exp != exp.Exp || last.Date.Format(isoDate) != exp.Date.Format(isoDate) {
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
