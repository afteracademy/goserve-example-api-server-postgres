package editor

import (
	userModel "github.com/afteracademy/goserve-example-api-server-postgres/api/user/model"
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
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		Controller: network.NewController("/blog/editor", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication(), c.Authorization(string(userModel.RoleCodeEditor)))
	group.GET("/id/:id", c.getBlogHandler)
	group.PUT("/publish/id/:id", c.publishBlogHandler)
	group.PUT("/unpublish/id/:id", c.unpublishBlogHandler)
	group.GET("/submitted", c.getSubmittedBlogsHandler)
	group.GET("/published", c.getPublishedBlogsHandler)
}

func (c *controller) getBlogHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blog, err := c.service.GetBlogById(uuidParam.ID)
	if err != nil {
		network.SendNotFoundError(ctx, uuidParam.ID.String()+" not found", err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", blog)
}

func (c *controller) publishBlogHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	err = c.service.BlogPublication(uuidParam.ID, true)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "blog published successfully")
}

func (c *controller) unpublishBlogHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	err = c.service.BlogPublication(uuidParam.ID, false)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "blog unpublished successfully")
}

func (c *controller) getSubmittedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blog, err := c.service.GetPaginatedSubmitted(pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blog)
}

func (c *controller) getPublishedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedPublished(pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blogs)
}
