package blogs

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/blogs/dto"
	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	network.Controller
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		Controller: network.NewController("/blogs", authMFunc, authorizeMFunc),
		service:    service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/latest", c.getLatestBlogsHandler)
	group.GET("/tag/:tag", c.getTaggedBlogsHandler)
	group.GET("/similar/id/:id", c.getSimilarBlogsHandler)

}

func (c *controller) getLatestBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedLatestBlogs(pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blogs)
}

func (c *controller) getTaggedBlogsHandler(ctx *gin.Context) {
	tag, err := network.ReqParams[dto.Tag](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedTaggedBlogs(tag.Tag, pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blogs)
}

func (c *controller) getSimilarBlogsHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blogs, err := c.service.GetSimilarBlogsDtoCache(uuidParam.ID)
	if err == nil {
		network.SendSuccessDataResponse(ctx, "success", &blogs)
		return
	}

	blogs, err = c.service.GetSimilarBlogs(uuidParam.ID)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blogs)
	c.service.SetSimilarBlogsDtoCache(uuidParam.ID, blogs)
}
