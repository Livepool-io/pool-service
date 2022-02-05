package db

import (
	"errors"
	"os"

	"github.com/livepool-io/pool-service/db/postgres"
	"github.com/livepool-io/pool-service/models"
)

var Database Store

type Store interface {
	GetJobs() ([]*models.Job, error)
	GetTranscoder(address string) (*models.Transcoder, error)
	GetTranscoders() ([]*models.Transcoder, error)
}

func Start() (Store, error) {
	if os.Getenv("POSTGRES") == "" {
		return nil, errors.New("Need to provide data source")
	}
	return postgres.Start()
}

func CacheDB() error {
	if Database != nil {
		return nil
	}

	var err error
	Database, err = Start()
	return err
}
