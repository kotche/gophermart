package postgres

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PGX struct {
	DB *sql.DB
}

func NewPGX(DSN string) (*PGX, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}
	pgx := &PGX{DB: db}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return pgx, nil
}

func (p *PGX) Ping() error {
	if err := p.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (p *PGX) Init() {
	_, err := p.DB.Exec(`CREATE TABLE IF NOT EXISTS public.users(
		    id SERIAL PRIMARY KEY,
    		login VARCHAR(255) NOT NULL UNIQUE,
    		password VARCHAR(255) NOT NULL);`)

	if err != nil {
		log.Fatal(err.Error())
	}
}
