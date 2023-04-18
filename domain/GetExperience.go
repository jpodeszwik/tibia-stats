package domain

import (
	"tibia-stats/repository"
)

func GetExperienceHistory(expRepository repository.ExpRepository, name string) ([]repository.ExpHistory, error) {
	return expRepository.GetExpHistory(name, 30)
}
