package repository

import (
	"testing"
	"time"
)

func Test_inMemoryExpRepository_GetExpStoreExp(t *testing.T) {
	i := NewInMemoryExpRepository()
	name := "Some Name"
	var exp int64 = 123456
	date := time.Now()

	err := i.StoreExp(name, date, exp)
	if nil != err {
		t.Errorf("Error received while storing data %v", err)
	}

	res, err := i.GetExp(name, date)
	if nil != err {
		t.Errorf("Error received while reading data %v", err)
	}
	if res != exp {
		t.Errorf("Invalid value received %v != %v", res, exp)
	}
}

func Test_inMemoryExpRepository_GetExp(t *testing.T) {
	i := NewInMemoryExpRepository()

	_, err := i.GetExp("Some Name", time.Now())
	if nil == err {
		t.Errorf("Error received while reading data %v", err)
	}
}
