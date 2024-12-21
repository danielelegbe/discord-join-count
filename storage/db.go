package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/danielelegbe/discord-join-count/sqlc"
)

func CreateAndMigrateStore(db *sql.DB, ddl string, ctx context.Context) *sqlc.Queries {

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		slog.Error(err.Error())
	}

	return sqlc.New(db)
}
