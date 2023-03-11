package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type postgresExpRepository struct {
	db *sql.DB
}

func (p *postgresExpRepository) StoreExperiences(expData []ExpData) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	valuesStr := make([]string, 0)
	values := make([]interface{}, 0)
	for i, ed := range expData {
		valuesStr = append(valuesStr, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		values = append(values, ed.Name)
		values = append(values, ed.Date.Format(isotime))
		values = append(values, ed.Exp)
	}

	query := "INSERT INTO exp(character_name, measure_date, exp_value) VALUES " + strings.Join(valuesStr, ",")
	_, err = tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return tx.Commit()
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
	err := createExpTable(db)
	if err != nil {
		log.Fatal(err)
	}
	return &postgresExpRepository{db: db}
}

func createExpTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS exp(id bigserial PRIMARY KEY, character_name text, measure_date date, exp_value bigint)")
	return err
}
