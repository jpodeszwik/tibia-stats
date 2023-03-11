package repository

import (
	"fmt"
	"time"
)

const isotime = "2006-01-02"

type dateName struct {
	date string
	name string
}

type inMemoryExpRepository struct {
	data map[dateName]int64
}

func (i *inMemoryExpRepository) StoreExperiences(expData []ExpData) error {
	for _, ed := range expData {
		err := i.StoreExp(ed.Name, ed.Date, ed.Exp)
		if nil != err {
			return err
		}
	}
	return nil
}

func (i *inMemoryExpRepository) StoreExp(name string, date time.Time, exp int64) error {
	i.data[dateName{date: date.Format(isotime), name: name}] = exp
	return nil
}

func (i *inMemoryExpRepository) GetExp(name string, time time.Time) (int64, error) {
	val, ok := i.data[dateName{date: time.Format(isotime), name: name}]
	if ok {
		return val, nil
	}

	return 0, fmt.Errorf("no value")
}

func NewInMemoryExpRepository() ExpRepository {
	return &inMemoryExpRepository{data: make(map[dateName]int64)}
}
