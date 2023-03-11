package repository

import "time"

type ExpData struct {
	Name string
	Date time.Time
	Exp  int64
}

type ExpRepository interface {
	StoreExperiences(expData []ExpData) error
	StoreExp(name string, date time.Time, exp int64) error
	GetExp(name string, time time.Time) (int64, error)
}
