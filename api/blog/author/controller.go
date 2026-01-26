package author

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/blog/dto"
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
		Controller:     network.NewController("/blog/author", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication(), c.Authorization(string(userModel.RoleCodeAuthor)))
	group.POST("/", c.postBlogHandler)
	group.PUT("/", c.updateBlogHandler)
	group.GET("/id/:id", c.getBlogHandler)
	group.DELETE("/id/:id", c.deleteBlogHandler)
	group.PUT("/submit/id/:id", c.submitBlogHandler)
	group.PUT("/withdraw/id/:id", c.withdrawBlogHandler)
	group.GET("/drafts", c.getDraftsBlogsHandler)
	group.GET("/submitted", c.getSubmittedBlogsHandler)
	group.GET("/published", c.getPublishedBlogsHandler)
}

func (c *controller) postBlogHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.BlogCreate](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	b, err := c.service.CreateBlog(body, user)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "blog created successfully", &b)
}

func (c *controller) updateBlogHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.BlogUpdate](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	b, err := c.service.UpdateBlog(body, user)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "blog updated successfully", &b)
}

func (c *controller) getBlogHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	blog, err := c.service.GetBlogById(uuidParam.ID, user)
	if err != nil {
		network.SendNotFoundError(ctx, uuidParam.ID.String()+" not found", err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", blog)
}

func (c *controller) submitBlogHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogSubmission(uuidParam.ID, user, true)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "blog submitted successfully")
}

func (c *controller) withdrawBlogHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogSubmission(uuidParam.ID, user, false)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "blog withdrawn successfully")
}

func (c *controller) deleteBlogHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.DeactivateBlog(uuidParam.ID, user)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "blog deleted successfully")
}

func (c *controller) getDraftsBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	blog, err := c.service.GetPaginatedDrafts(user, pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blog)
}

func (c *controller) getSubmittedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	blog, err := c.service.GetPaginatedSubmitted(user, pagination)
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

	user := c.MustGetUser(ctx)

	blogs, err := c.service.GetPaginatedPublished(user, pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blogs)
}
