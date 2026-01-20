package middleware

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/user/model"
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type authorizationProvider struct {
	common.ContextPayload
}

func NewAuthorizationProvider() network.AuthorizationProvider {
	return &authorizationProvider{
		ContextPayload: common.NewContextPayload(),
	}
}

func (m *authorizationProvider) Middleware(roleNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(roleNames) == 0 {
			network.SendForbiddenError(ctx, "permission denied: role missing", nil)
			return
		}

		user := m.MustGetUser(ctx)

		hasRole := false
		for _, code := range roleNames {
			for _, role := range user.Roles {
				if role.Code == model.RoleCode(code) {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			network.SendForbiddenError(ctx, "permission denied: does not have suffient role", nil)
			return
		}

		ctx.Next()
	}
}
