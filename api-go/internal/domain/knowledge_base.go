package domain

import "time"

type KnowledgeBaseArticle struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateKnowledgeBaseArticleResponse struct {
	Message string               `json:"message"`
	Article KnowledgeBaseArticle `json:"article"`
}

type UpdateKnowledgeBaseArticleResponse struct {
	Message string               `json:"message"`
	Article KnowledgeBaseArticle `json:"article"`
}

type DeleteKnowledgeBaseArticleResponse struct {
	Message string `json:"message"`
}
