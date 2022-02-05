package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/livepool-io/pool-service/models"
)

type Postgres struct {
	*sql.DB
}

func Start() (*Postgres, error) {
	datasource := os.Getenv("POSTGRES")
	sql, err := sql.Open("postgres", datasource)
	if err != nil {
		return nil, err
	}
	db := &Postgres{sql}

	if err := db.ensureDatabase(); err != nil {
		return nil, err
	}

	if err := db.ensureTables(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Postgres) ensureDatabase() error {
	_, err := db.Exec(`CREATE DATABASE livepool`)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		db.Close()
		return err
	}
	return nil
}

func (db *Postgres) ensureTables() error {
	_, err := db.Exec(
		fmt.Sprint(`
		CREATE TABLE transcoders (
			id SERIAL PRIMARY KEY,
			transcoder JSONB
		);

		CREATE TABLE jobs (
			id SERIAL PRIMARY KEY,
			job JSONB
		)
		`),
	)

	if err != nil && !strings.Contains(err.Error(), "already exists") {
		db.Close()
		return err
	}

	return nil
}

func (db *Postgres) GetTranscoder(address string) (*models.Transcoder, error) {
	var t *models.Transcoder
	return t, nil
}

func (db *Postgres) GetTranscoders() ([]*models.Transcoder, error) {
	return []*models.Transcoder{}, nil
}

func (db *Postgres) GetJobs() ([]*models.Job, error) {
	return []*models.Job{}, nil
}
