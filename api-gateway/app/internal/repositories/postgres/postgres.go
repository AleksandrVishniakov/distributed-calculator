package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DBConfigs struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	DBName   string `yaml:"db-name"`
	SSLMode  string `yaml:"ssl-mode"`

	Password string
}

func NewPostgresDB(cfg *DBConfigs) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
