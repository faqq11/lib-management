package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDb() (*sqlx.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("env not found")
	}

  dsn := os.Getenv("DATABASE_URL")
  db, err := sqlx.Connect("postgres", dsn)
  if err != nil {
    return nil, fmt.Errorf("connect db: %w", err)
  }
  return db, nil
}
