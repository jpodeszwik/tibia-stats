package tracker

import (
	"log"
	"sync"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
)

const isotime = "2006-01-02"

type GuildExp struct {
	lastExp    map[string]*domain.GuildExp
	m          sync.RWMutex
	repository *dynamo.GuildExpRepository
}

func (ge *GuildExp) Handle(exp domain.GuildExp) {
	ge.m.RLock()
	last, exists := ge.lastExp[exp.GuildName]
	ge.m.RUnlock()

	if !exists || last.Exp != exp.Exp || last.Date.Format(isotime) != exp.Date.Format(isotime) {
		log.Printf("Storing %v", exp)
		ge.repository.StoreGuildExp(exp)
		ge.m.Lock()
		ge.lastExp[exp.GuildName] = &exp
		ge.m.Unlock()
	}
}

func NewGuildExp(repository *dynamo.GuildExpRepository) *GuildExp {
	return &GuildExp{
		lastExp:    make(map[string]*domain.GuildExp),
		m:          sync.RWMutex{},
		repository: repository,
	}
}
