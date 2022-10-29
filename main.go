package main

import (
	"database/sql"
	"github.com/jpodeszwik/tibia-exp-tracker/repository"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func main() {
	connStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	postgresRepository := repository.NewPostgresExpRepository(db)

	now := time.Now()
	name := "Test Name"

	err = postgresRepository.StoreExp(name, now, 34567)
	if nil != err {
		log.Fatal("Error", err)
	}
	exp, err := postgresRepository.GetExp(name, now)
	if nil != err {
		log.Fatal("Error", err)
	}

	log.Printf("Exp %v", exp)
}
