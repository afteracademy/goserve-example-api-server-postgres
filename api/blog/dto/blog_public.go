package dto

import (
	"time"

	"github.com/afteracademy/goserve-example-api-server-postgres/api/blog/model"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/user/dto"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/google/uuid"
)

type BlogPublic struct {
	ID          uuid.UUID       `json:"id" binding:"required" validate:"required"`
	Title       string          `json:"title" validate:"required,min=3,max=500"`
	Description string          `json:"description" validate:"required,min=3,max=2000"`
	Text        string          `json:"text" validate:"required,max=50000"`
	Slug        string          `json:"slug" validate:"required,min=3,max=200"`
	Author      *dto.UserPublic `json:"author,omitempty" validate:"required,omitempty"`
	ImgURL      *string         `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       *float64        `json:"score,omitempty" validate:"omitempty,min=0,max=1"`
	Tags        *[]string       `json:"tags,omitempty" validate:"omitempty,dive,uppercase"`
	PublishedAt *time.Time      `json:"publishedAt,omitempty"`
}

func NewBlogPublic(blog *model.Blog, author *dto.UserPublic) (*BlogPublic, error) {
	b, err := utility.MapTo[BlogPublic](blog)
	if err != nil {
		return nil, err
	}

	b.Author = author

	return b, err
}
