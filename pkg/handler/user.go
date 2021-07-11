package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
)

type RoleCodeFilter struct {
	RoleCode string `json:"role_code" form:"role_code" binding:"required"`
}

type usersResponse struct {
	Data []model.User `json:"data"`
}

// GetUserHandler get user handler
// @Summary get user handler
// @Description get user
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param role_code query RoleCodeFilter true "role code"
// @Security ApiKeyAuth
// @Success 200 {object} handler.usersResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /users [get]
func (h *Handler) GetUserHandler(c *gin.Context) {
	filter := RoleCodeFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.HandleError(c, err)
		return
	}
	if filter.RoleCode == "" {
		users, err := h.repo.User.GetAll(h.store)
		if err != nil {
			h.log.Error("[handler.GetUserHandler] User.GetAll()", err)
			util.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, usersResponse{Data: users})

	} else {
		users, err := h.repo.User.GetByRoleCode(h.store, filter.RoleCode)
		if err != nil {
			h.log.Error("[handler.GetUserHandler] User.GetByRoleCode()", err)
			util.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, usersResponse{Data: users})
	}
}
