package contact

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/contact/dto"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/gin-gonic/gin"
)

type controller struct {
	network.Controller
	service Service
}

func NewController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		Controller: network.NewController("/contact", authProvider, authorizeProvider),
		service:    service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/", c.createMessageHandler)
}

func (c *controller) createMessageHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.MessageCreate](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	msg, err := c.service.CreateMessage(body)
	if err != nil {
		network.SendInternalServerError(ctx, "something went wrong", err)
		return
	}

	data, err := utility.MapTo[dto.Message](msg)
	if err != nil {
		network.SendInternalServerError(ctx, "something went wrong", err)
		return
	}

	network.SendSuccessDataResponse(ctx, "message received successfully!", data)
}
