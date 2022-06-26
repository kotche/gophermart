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
    		password VARCHAR(255) NOT NULL);

			CREATE TABLE IF NOT EXISTS public.accruals(
				order_num BIGINT PRIMARY KEY,
				user_id INT NOT NULL,
				status VARCHAR(30) NOT NULL DEFAULT 'NEW',
				amount INT,
				uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				CONSTRAINT accruals_uniq_order_user UNIQUE (order_num, user_id),
				FOREIGN KEY (user_id) REFERENCES public.users (id));

			CREATE TABLE IF NOT EXISTS public.withdrawals(
			    order_num BIGINT PRIMARY KEY,
				user_id INT NOT NULL,
				amount INT,
				processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			    CONSTRAINT withdrawals_uniq_order_user UNIQUE (order_num, user_id),
				FOREIGN KEY (user_id) REFERENCES public.users (id));
`)

	if err != nil {
		log.Fatal(err.Error())
	}
}
