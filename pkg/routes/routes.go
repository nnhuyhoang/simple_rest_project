package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/config"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/handler"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/mw"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

func NewRoutes(r *gin.Engine, h *handler.Handler, cfg config.Config, s repo.DBRepo) {
	authMw := mw.NewAuthMiddleware(cfg, s)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", h.LoginHandler)
		authGroup.POST("/signup", h.SignupHandler)
	}

	userGroup := r.Group("/users")
	{
		userGroup.GET("", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnicianManager), h.GetUserHandler)

	}

	siteGroup := r.Group("/sites")
	{
		siteGroup.GET("", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnicianManager), h.GetSiteByDateHandler)
		siteGroup.GET("/:site_id", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnicianManager), h.GetSiteByIdHandler)
	}

	inspectionGroup := r.Group("/inspections")
	{
		inspectionGroup.GET("", authMw.WithAuth, h.GetInspectionByUserIdHandler)
		inspectionGroup.GET("/:inspection_id", authMw.WithAuth, h.GetInspectionByIdHandler)
		inspectionGroup.POST("", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnicianManager), h.CreateInspectionHandler)
		inspectionGroup.PUT("/:inspection_id", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnician), h.UpdateInspectionHandler)
	}

	issueGroup := r.Group("/issues")
	{
		issueGroup.GET("/:issue_id", authMw.WithAuth, h.GetIssueByIdhandler)
		issueGroup.POST("", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnician), h.CreateIssueHandler)
		issueGroup.PUT("/:issue_id", authMw.WithAuth, h.UpdateIssueHandler)
		issueGroup.DELETE("/:issue_id", authMw.WithAuth, h.DeleteIssueByIdHandler)
	}

	sparePartGroup := r.Group("/spare-parts")
	{
		sparePartGroup.GET("", authMw.WithAuth, h.GetAllSparePartHandler)
	}

	issueSparePartGroup := r.Group("/issue-spare-parts")
	{
		issueSparePartGroup.GET("", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnicianManager), h.GetAllIssueSparePartHandler)
	}
	purchaseRequestGroup := r.Group("/purchase-requests")
	{
		purchaseRequestGroup.GET("", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnicianManager), h.GetAllPurchaseRequestHandler)
		purchaseRequestGroup.POST("", authMw.WithAuth, authMw.AuthPerm(consts.FieldTechnicianManager), h.CreatePurchaseRequestHandler)
	}
}
