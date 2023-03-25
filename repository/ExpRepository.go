package repository

type ExpRepository interface {
	StoreExperiences(expData []ExpData) error
	GetExpHistory(name string, limit int) ([]ExpHistory, error)
}
