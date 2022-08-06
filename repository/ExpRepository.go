package repository

import "time"

type ExpRepository interface {
	StoreExp(name string, date time.Time, exp int64) error
	GetExp(name string, time time.Time) (int64, error)
}
