package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GDB struct {
	DB *gorm.DB
}

func NewGDB(config Config) (*GDB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		CreateBatchSize: 1000,
	})

	if err != nil {
		return nil, err
	}

	return &GDB{DB: db}, nil
}
