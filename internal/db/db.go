package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/blockseeker999th/SpyCat/internal/config"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
)

func ConnectDB(config *config.Config) (*sql.DB, error) {
	const op = "db.ConnectDB"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Err(err)
		log.Fatalf("error path: %s, error: %v", op, err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging PostgreSQL DB, path: %s, error:%v", op, err)
		return nil, err
	}

	return db, nil
}
