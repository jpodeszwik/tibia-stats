package repository

import (
	"database/sql"
	"log"
	"time"
)

type postgresExpRepository struct {
	db *sql.DB
}

func (p *postgresExpRepository) StoreExp(name string, date time.Time, exp int64) error {
	query := "INSERT INTO exp(character_name, measure_date, exp_value) VALUES ($1, $2, $3)"
	_, err := p.db.Exec(query, name, date.Format(isotime), exp)
	return err
}

func (p *postgresExpRepository) GetExp(name string, date time.Time) (int64, error) {
	query := "SELECT MAX(exp_value) FROM exp WHERE character_name = $1 AND measure_date = $2"
	row := p.db.QueryRow(query, name, date.Format(isotime))
	var value int64 = 0
	err := row.Scan(&value)
	return value, err
}

func NewPostgresExpRepository(db *sql.DB) ExpRepository {
	err := initDb(db)
	if err != nil {
		log.Fatal(err)
	}
	return &postgresExpRepository{db: db}
}

func initDb(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS exp(id bigserial PRIMARY KEY, character_name text, measure_date date, exp_value bigint)")
	return err
}
