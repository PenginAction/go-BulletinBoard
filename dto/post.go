package dto

import "time"

type CreatePostRequest struct {
	UserID uint   `json:"user_id" validate:"required"`
	Text   string `json:"text" validate:"required,min=1"`
}

type AllPostsRequest struct {
	UserID uint  `json:"user_id" validate:"required"`
	Limit  int32 `json:"limit" validate:"required"`
	Offset int32 `json:"offset" validate:"required"`
}

type UpdatePostRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
