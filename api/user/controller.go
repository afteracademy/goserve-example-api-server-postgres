package user

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	coredto "github.com/afteracademy/goserve/v2/dto"
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
		Controller:     network.NewController("/profile", authProvider, authorizeProvider),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getPublicProfileHandler)
	private := group.Use(c.Authentication())
	private.GET("/mine", c.getPrivateProfileHandler)
}

func (c *controller) getPublicProfileHandler(ctx *gin.Context) {
	dto, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	data, err := c.service.FetchUserPublicProfile(dto.ID)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", data)
}

func (c *controller) getPrivateProfileHandler(ctx *gin.Context) {
	user := c.MustGetUser(ctx)

	data, err := c.service.FetchUserPrivateProfile(user)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", data)
}
