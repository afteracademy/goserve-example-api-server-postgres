package auth

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/auth/dto"
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	"github.com/afteracademy/goserve-example-api-server-postgres/utils"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	network.Controller
	common.ContextPayload
	service Service
}

func NewController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		Controller: network.NewController("/auth", authProvider, authorizeProvider),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/signup/basic", c.signUpBasicHandler)
	group.POST("/signin/basic", c.signInBasicHandler)
	group.POST("/token/refresh", c.tokenRefreshHandler)
	group.DELETE("/signout", c.Authentication(), c.signOutBasic)
}

func (c *controller) signUpBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.SignUpBasic](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	data, err := c.service.SignUpBasic(body)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", data)
}

func (c *controller) signInBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.SignInBasic](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	dto, err := c.service.SignInBasic(body)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", dto)
}

func (c *controller) signOutBasic(ctx *gin.Context) {
	keystore := c.MustGetKeystore(ctx)

	err := c.service.SignOut(keystore)
	if err != nil {
		network.SendInternalServerError(ctx, "something went wrong", err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "signout success")
}

func (c *controller) tokenRefreshHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.TokenRefresh](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	authHeader := ctx.GetHeader(network.AuthorizationHeader)
	accessToken := utils.ExtractBearerToken(authHeader)

	dto, err := c.service.RenewToken(body, accessToken)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", dto)
}
