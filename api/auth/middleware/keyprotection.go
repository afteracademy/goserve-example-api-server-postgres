package middleware

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/auth"
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type keyProtection struct {
	network.ResponseSender
	common.ContextPayload
	authService auth.Service
}

func NewKeyProtection(authService auth.Service) network.RootMiddleware {
	return &keyProtection{
		ResponseSender: network.NewResponseSender(),
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
	}
}

func (m *keyProtection) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *keyProtection) Handler(ctx *gin.Context) {
	key := ctx.GetHeader(network.ApiKeyHeader)
	if len(key) == 0 {
		m.Send(ctx).UnauthorizedError("permission denied: missing x-api-key header", nil)
		return
	}

	apikey, err := m.authService.FetchApiKey(key)
	if err != nil {
		m.Send(ctx).ForbiddenError("permission denied: invalid x-api-key", err)
		return
	}

	m.SetApiKey(ctx, apikey)

	ctx.Next()
}
