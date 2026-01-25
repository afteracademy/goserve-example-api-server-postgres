package health

import (
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	network.Controller
	service Service
}

func NewController(
	service Service,
) network.Controller {
	return &controller{
		Controller: network.NewController("/health", nil, nil),
		service:    service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/", c.getHealthHandler)
}

func (c *controller) getHealthHandler(ctx *gin.Context) {
	health, err := c.service.CheckHealth()
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", health)
}
