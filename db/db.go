package db

import (
	"errors"
	"os"

	"github.com/livepool-io/pool-service/db/postgres"
	"github.com/livepool-io/pool-service/models"
)

var Database Store

type Store interface {
	// Transcoders
	GetTranscoder(address string) (*models.Transcoder, error)
	GetTranscoders() ([]*models.Transcoder, error)

	// Nodes
	AddNode(n *models.Node) error
	GetNodes(transcoder string, region string) ([]*models.Node, error)

	// Jobs
	CreateJob(job *models.Job) error
	GetJobs(transcoder, node string, from, to int64, isAuth bool) ([]*models.Job, error)
}

func Start() (Store, error) {
	if os.Getenv("POSTGRES") == "" {
		return nil, errors.New("need to provide data source")
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
