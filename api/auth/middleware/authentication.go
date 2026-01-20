package middleware

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/auth"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/user"
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	"github.com/afteracademy/goserve-example-api-server-postgres/utils"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authenticationProvider struct {
	common.ContextPayload
	authService auth.Service
	userService user.Service
}

func NewAuthenticationProvider(authService auth.Service, userService user.Service) network.AuthenticationProvider {
	return &authenticationProvider{
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
		userService:    userService,
	}
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(network.AuthorizationHeader)
		if len(authHeader) == 0 {
			network.SendUnauthorizedError(ctx, "permission denied: missing Authorization", nil)
			return
		}

		token := utils.ExtractBearerToken(authHeader)
		if token == "" {
			network.SendUnauthorizedError(ctx, "permission denied: invalid Authorization", nil)
			return
		}

		claims, err := m.authService.VerifyToken(token)
		if err != nil {
			network.SendUnauthorizedError(ctx, err.Error(), err)
			return
		}

		valid := m.authService.ValidateClaims(claims)
		if !valid {
			network.SendUnauthorizedError(ctx, "permission denied: invalid claims", nil)
			return
		}

		userId, err := uuid.Parse(claims.Subject)
		if err != nil {
			network.SendUnauthorizedError(ctx, "permission denied: invalid claims subject", nil)
			return
		}

		user, err := m.userService.FetchUserById(userId)
		if err != nil {
			network.SendUnauthorizedError(ctx, "permission denied: claims subject does not exists", err)
			return
		}

		keystore, err := m.authService.FetchKeystore(user, claims.ID)
		if err != nil || keystore == nil {
			network.SendUnauthorizedError(ctx, "permission denied: invalid access token", err)
			return
		}

		m.SetUser(ctx, user)
		m.SetKeystore(ctx, keystore)

		ctx.Next()
	}
}
