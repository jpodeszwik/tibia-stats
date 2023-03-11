package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type postgresGuildMemberRepository struct {
	db *sql.DB
}

func (p *postgresGuildMemberRepository) StoreGuildMembers(guild string, members []string) error {
	valuesStr := make([]string, 0)
	values := make([]interface{}, 0)
	for i, member := range members {
		valuesStr = append(valuesStr, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		values = append(values, guild)
		values = append(values, time.Now().Format(isotime))
		values = append(values, member)
	}

	query := "INSERT INTO guild_members(guild_name, measure_date, member_name) VALUES " + strings.Join(valuesStr, ",")
	_, err := p.db.Exec(query, values...)
	return err
}

func (p *postgresGuildMemberRepository) StoreExperiences(expData []ExpData) error {
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

func NewPostgresGuildMemberRepository(db *sql.DB) GuildMemberRepository {
	err := createGuildMemberTable(db)
	if err != nil {
		log.Fatal(err)
	}
	return &postgresGuildMemberRepository{db: db}
}

func createGuildMemberTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS guild_members(id bigserial PRIMARY KEY, guild_name text, measure_date date, member_name text)")
	return err
}
