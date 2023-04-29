package tracker

import (
	"sync"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
)

type GuildExp struct {
	lastExp    map[string]int64
	m          sync.RWMutex
	repository *dynamo.GuildExpRepository
}

func (ge *GuildExp) Handle(exp domain.GuildExp) {
	ge.m.RLock()
	last, exists := ge.lastExp[exp.GuildName]
	ge.m.RUnlock()

	if !exists || last != exp.Exp {
		ge.m.Lock()
		ge.lastExp[exp.GuildName] = exp.Exp
		ge.m.Unlock()
		ge.repository.StoreGuildExp(exp)
	}
}

func NewGuildExp(repository *dynamo.GuildExpRepository) *GuildExp {
	return &GuildExp{
		lastExp:    make(map[string]int64),
		m:          sync.RWMutex{},
		repository: repository,
	}
}
