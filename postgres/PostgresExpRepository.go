package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tibia-exp-tracker/repository"
)

const isotime = "2006-01-02"

type postgresExpRepository struct {
	db *sql.DB
}

func (p *postgresExpRepository) GetExpHistory(name string, limit int) ([]repository.ExpHistory, error) {
	query := "SELECT measure_date, exp_value FROM exp WHERE character_name = $1 order by measure_date desc limit $2"
	rows, err := p.db.Query(query, name, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	historySlice := make([]repository.ExpHistory, 0)

	for rows.Next() {
		var history repository.ExpHistory
		err = rows.Scan(&history.Date, &history.Exp)
		if err != nil {
			return historySlice, err
		}
		historySlice = append(historySlice, history)
	}

	return historySlice, rows.Err()
}

func (p *postgresExpRepository) StoreExperiences(expData []repository.ExpData) error {
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

func NewPostgresExpRepository(db *sql.DB) repository.ExpRepository {
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
