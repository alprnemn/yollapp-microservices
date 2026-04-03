package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func Init() (*sql.DB, error) {

	DB, err := sql.Open(
		"sqlite",
		"file:./services/auth/internal/db/auth.db",
	)
	if err != nil {
		return nil, fmt.Errorf("error open db: %s", err.Error())
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping error: %s", err.Error())
	}

	query := `
	CREATE TABLE IF NOT EXISTS user_invitations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		token TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

	_, err = DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing init query: %s", err.Error())
	}

	log.Printf("\033[38;5;226m DB Started succesfully \033[0m")
	return DB, err
}
