package dto

import "time"

type CreatePostRequest struct {
	UserID uint   `json:"user_id" validate:"required"`
	Text   string `json:"text" validate:"required,min=1"`
}

type AllPostsRequest struct {
	PageID   int32 `form:"page_id" validate:"required,min=1,max=10"`
	PageSize int32 `form:"page_size" validate:"required,min=1"`
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
