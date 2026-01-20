package dto

import (
	"time"

	"github.com/afteracademy/goserve-example-api-server-postgres/api/blog/model"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/user/dto"
	userModel "github.com/afteracademy/goserve-example-api-server-postgres/api/user/model"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/google/uuid"
)

type BlogPrivate struct {
	ID          uuid.UUID       `json:"id" binding:"required" validate:"required"`
	Title       string          `json:"title" validate:"required,min=3,max=500"`
	Description string          `json:"description" validate:"required,min=3,max=2000"`
	Text        *string         `json:"text,omitempty" validate:"omitempty,max=50000"`
	DraftText   string          `json:"draftText" validate:"required"`
	Slug        string          `json:"slug" validate:"required,min=3,max=200"`
	Author      *dto.UserPublic `json:"author,omitempty" validate:"required,omitempty"`
	ImgURL      *string         `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       *float64        `json:"score,omitempty" validate:"omitempty,min=0,max=1"`
	Tags        *[]string       `json:"tags,omitempty" validate:"omitempty,dive,uppercase"`
	Submitted   bool            `json:"submitted" validate:"required"`
	Drafted     bool            `json:"drafted" validate:"required"`
	Published   bool            `json:"published" validate:"required"`
	PublishedAt *time.Time      `json:"publishedAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time       `json:"updatedAt" validate:"required"`
}

func NewBlogPrivate(blog *model.Blog, author *userModel.User) (*BlogPrivate, error) {
	b, err := utility.MapTo[BlogPrivate](blog)
	if err != nil {
		return nil, err
	}

	b.Author, err = utility.MapTo[dto.UserPublic](author)
	if err != nil {
		return nil, err
	}

	return b, err
}
