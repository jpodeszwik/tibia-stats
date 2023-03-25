package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tibia-exp-tracker/repository"
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

func NewPostgresGuildMemberRepository(db *sql.DB) repository.GuildMemberRepository {
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
