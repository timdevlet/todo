package logs

import (
	log "github.com/sirupsen/logrus"
	"github.com/timdevlet/todo/pkg/postgres"
)

type ILogRepository interface {
	InsertLog(l Log) (string, error)
}

type LogRepository struct {
	postgres *postgres.PDB
}

func NewLogRepository(db *postgres.PDB) ILogRepository {
	return &LogRepository{
		postgres: db,
	}
}

// ----------------------------

func (repo *LogRepository) InsertLog(l Log) (string, error) {
	// @todo move to sql file
	sqlStatement := `
		INSERT INTO logs (payload)
		VALUES ($1)
		RETURNING uuid
	`

	uuid := ""
	err := repo.postgres.DB.QueryRow(sqlStatement, l.ToJson()).Scan(&uuid)

	log.WithField("uuid", uuid).Debug("New log")

	return uuid, err
}
