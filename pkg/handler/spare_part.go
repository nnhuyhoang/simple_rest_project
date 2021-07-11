package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
)

type sparePartListReponse struct {
	Data []model.SparePart `json:"data"`
}

// GetAllSparePartHandler get all spare part handler
// @Summary get all spare part handler
// @Description get list of spar part
// @Tags SparePart
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Security ApiKeyAuth
// @Success 200 {object} sparePartListReponse	"ok"
// @Failure 500 {object} model.Error
// @Router /spare-parts [get]
func (h *Handler) GetAllSparePartHandler(c *gin.Context) {
	spareParts, err := h.repo.SparePart.GetAll(h.store)
	if err != nil {
		h.log.Error("[handler.GetAllSparePart] SparePart.GetAll()", err)
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, sparePartListReponse{Data: spareParts})
}
