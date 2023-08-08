package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PDB struct {
	DB *sql.DB
}

func NewPDB(config Config) (*PDB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PDB{DB: db}, nil
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}
