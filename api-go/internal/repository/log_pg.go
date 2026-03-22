package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"med-vito/api-go/internal/domain"
)

// LogPG — доступ к таблице "Log" (имена в кавычках как в миграции Prisma).
type LogPG struct {
	pool *pgxpool.Pool
}

func NewLogPG(pool *pgxpool.Pool) *LogPG {
	return &LogPG{pool: pool}
}

func (r *LogPG) FindAll(ctx context.Context) ([]domain.Log, error) {
	const q = `SELECT "id", "userId", "action" FROM "Log" ORDER BY "id" ASC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("log findAll query: %w", err)
	}
	defer rows.Close()

	var out []domain.Log
	for rows.Next() {
		var row domain.Log
		if err := rows.Scan(&row.ID, &row.UserID, &row.Action); err != nil {
			return nil, fmt.Errorf("log scan: %w", err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
