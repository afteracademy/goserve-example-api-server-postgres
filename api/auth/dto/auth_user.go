package dto

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/user/dto"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/user/model"
)

type UserAuth struct {
	User   *dto.UserPrivate `json:"user" validate:"required"`
	Tokens *Tokens          `json:"tokens" validate:"required"`
}

func NewUserAuth(user *model.User, tokens *Tokens) *UserAuth {
	return &UserAuth{
		User:   dto.NewUserPrivate(user),
		Tokens: tokens,
	}
}
