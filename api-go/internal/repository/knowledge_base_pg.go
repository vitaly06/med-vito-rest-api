package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"med-vito/api-go/internal/domain"
)

type KnowledgeBasePG struct {
	pool *pgxpool.Pool
}

func NewKnowledgeBasePG(pool *pgxpool.Pool) *KnowledgeBasePG {
	return &KnowledgeBasePG{pool: pool}
}

func (r *KnowledgeBasePG) CreateArticle(ctx context.Context, title, content string) (*domain.KnowledgeBaseArticle, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO "KnowledgeBaseArticle" ("title", "content", "createdAt", "updatedAt")
		VALUES ($1, $2, NOW(), NOW())
		RETURNING "id", "title", "content", "createdAt", "updatedAt"`,
		title, content,
	)
	var article domain.KnowledgeBaseArticle
	if err := row.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *KnowledgeBasePG) ListArticles(ctx context.Context) ([]domain.KnowledgeBaseArticle, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT "id", "title", "content", "createdAt", "updatedAt"
		FROM "KnowledgeBaseArticle"
		ORDER BY "createdAt" DESC, "id" DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.KnowledgeBaseArticle
	for rows.Next() {
		var article domain.KnowledgeBaseArticle
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, article)
	}
	return out, rows.Err()
}

func (r *KnowledgeBasePG) GetArticleByID(ctx context.Context, id int32) (*domain.KnowledgeBaseArticle, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT "id", "title", "content", "createdAt", "updatedAt"
		FROM "KnowledgeBaseArticle"
		WHERE "id" = $1`,
		id,
	)
	var article domain.KnowledgeBaseArticle
	if err := row.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (r *KnowledgeBasePG) UpdateArticle(ctx context.Context, id int32, title, content string) (*domain.KnowledgeBaseArticle, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE "KnowledgeBaseArticle"
		SET "title" = $2, "content" = $3, "updatedAt" = NOW()
		WHERE "id" = $1
		RETURNING "id", "title", "content", "createdAt", "updatedAt"`,
		id, title, content,
	)
	var article domain.KnowledgeBaseArticle
	if err := row.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (r *KnowledgeBasePG) DeleteArticle(ctx context.Context, id int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM "KnowledgeBaseArticle" WHERE "id" = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
