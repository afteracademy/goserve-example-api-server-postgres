package dto

import (
	"time"

	"github.com/afteracademy/goserve-example-api-server-postgres/api/blog/model"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/google/uuid"
)

type BlogItem struct {
	ID          uuid.UUID  `json:"id" binding:"required" validate:"required"`
	Title       string     `json:"title" validate:"required,min=3,max=500"`
	Description string     `json:"description" validate:"required,min=3,max=2000"`
	Slug        string     `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      *string    `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       float64    `json:"score," validate:"required,min=0,max=1"`
	Tags        []string   `json:"tags" validate:"required,dive,uppercase"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
}

func NewBlogItem(blog *model.Blog) (*BlogItem, error) {
	return utility.MapTo[BlogItem](blog)
}
