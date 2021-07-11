package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
)

type issueSparePartListReponse struct {
	Data []model.IssueSparePart `json:"data"`
}

// GetAllIssueSparePartHandler get all issue spare part handler
// @Summary get all issue spare part handler
// @Description before creating purchase request, manger call this to get list of spare parts which is ordered but not enough in iventory
// @Tags IssueSparePart
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Security ApiKeyAuth
// @Success 200 {object} issueSparePartListReponse	"ok"
// @Failure 500 {object} model.Error
// @Router /issue-spare-parts [get]
func (h *Handler) GetAllIssueSparePartHandler(c *gin.Context) {
	iSpareParts, err := h.repo.IssueSparePart.GetByStatus(h.store, consts.StatusInProgress)
	if err != nil {
		h.log.Error("[handler.GetAllIssueSparePartHandler] IssueSparePart.GetByStatus()", err)
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, issueSparePartListReponse{Data: iSpareParts})
}
