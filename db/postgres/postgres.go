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

func (db *Postgres) GetJobs(transcoder, node string, from, to int64, isAuth bool) ([]*models.Job, error) {
	qry := fmt.Sprintf(`SELECT job FROM jobs`)
	if transcoder != "" {
		qry = fmt.Sprintf(`%s WHERE job->>'transcoder' = '%v'`, qry, transcoder)
	}
	// if node is defined require transcoder to be defined
	if transcoder != "" && node != "" {
		qry = fmt.Sprintf(`%s AND jobs->>'node' = '%v'`, qry, node)
	}
	// timestamp filter
	if transcoder != "" {
		qry = fmt.Sprintf(`%s AND job->>'timestamp' >= '%v' AND job->>'timestamp' <= '%v' ORDER BY job->>'timestamp' DESC`, qry, from, to)
	} else {
		qry = fmt.Sprintf(`%s WHERE job->>'timestamp' >= '%v' AND job->>'timestamp' <= '%v' ORDER BY job->>'timestamp' DESC`, qry, from, to)
	}

	rows, err := db.Query(qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []*models.Job{}
	for rows.Next() {
		var job *models.Job
		if err := rows.Scan(job); err != nil {
			return nil, err
		}
		if !isAuth {
			job.Node = ""
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (db *Postgres) CreateJob(job *models.Job) error {
	qry := fmt.Sprintf("INSERT INTO jobs (job) VALUES($1)")
	_, err := db.Exec(qry, job)
	return err
}
