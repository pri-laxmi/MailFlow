package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pri-laxmi/MailFlow/configs"
)

// ConnectDB connects to PostgreSQL using environment variables.
// Expected env vars: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
func ConnectDB(cfg *configs.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return db, nil
}
