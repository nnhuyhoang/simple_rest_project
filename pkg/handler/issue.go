package handler

import (
	"encoding/json"
	inerr "errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
)

type createIssueRequest struct {
	SiteId            int                `json:"siteId" form:"siteId" binding:"required" example:"1"`
	InspectionId      int                `json:"inspectionId" form:"inspectionId" binding:"required" example:"2"`
	Description       string             `json:"description" form:"description" example:"something wrong"`
	SparePartRequests []sparePartRequest `json:"sparePartRequests" form:"sparePartRequests"`
}

type updateIssueRequest struct {
	Description string `json:"description" form:"description" example:"something wrong"`
	Status      string `json:"status" form:"status" example:"completed"`
}

type sparePartRequest struct {
	SparePartId int `json:"sparePartId" form:"sparePartId" binding:"required" example:"1"`
	Quantity    int `json:"quantity" form:"quantity" binding:"required" example:"100"`
}

type issueResponse struct {
	Data model.Issue `json:"data"`
}

// GetIssueByIdhandler get issue by Id handler
// @Summary get issue by Id handler
// @Description get issue by id
// @Tags Issue
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param issue_id path string true "issue_id"
// @Security ApiKeyAuth
// @Success 200 {object} issueResponse	"ok"
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /issues/{issue_id} [get]
func (h *Handler) GetIssueByIdhandler(c *gin.Context) {
	issueIdStr := c.Param("issue_id")
	issueId, err := strconv.Atoi(issueIdStr)
	if err != nil {
		h.log.Error("[handler.GetIssueByIdhandler] strconv.Atoi(issueIdStr)", err)
		util.HandleError(c, errors.ErrInvalidIdFormat)
		return
	}

	issue, err := h.repo.Issue.GetByIdWithSparePart(h.store, issueId)
	if err != nil {
		h.log.Error("[handler.GetIssueByIdhandler] GetByIdWithSparePart()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			util.HandleError(c, errors.ErrIssueNotFound)
			return
		}
		util.HandleError(c, errors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, issueResponse{Data: *issue})
}

// CreateIssueHandler create issue handler
// @Summary create issue handler
// @Description FT create issue
// @Tags Issue
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param body body createIssueRequest true "body"
// @Security ApiKeyAuth
// @Success 200 {object} issueResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /issues [post]
func (h *Handler) CreateIssueHandler(c *gin.Context) {

	authInfo, err := util.GetAuthenInfoFromContext(c)
	if err != nil {
		h.log.Error("[handler.CreateIssueHandler] GetAuthenInfoFromContext()", err)
		util.HandleError(c, err)
		return
	}

	request := createIssueRequest{}
	if err = c.ShouldBindJSON(&request); err != nil {
		util.HandleError(c, err)
		return
	}

	var issueSparePartParams []model.IssueSparePart
	for _, v := range request.SparePartRequests {
		newV := model.IssueSparePart{
			SparePartId: v.SparePartId,
			Quantity:    v.Quantity,
		}
		issueSparePartParams = append(issueSparePartParams, newV)
	}
	IssueParam := model.Issue{
		UserId:          authInfo.UserID,
		SiteId:          request.SiteId,
		InspectionId:    request.InspectionId,
		IssueSpareParts: issueSparePartParams,
		Description:     request.Description,
		Status:          consts.StatusInProgress,
	}

	issue, err := h.doCreateIssue(&IssueParam)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, issueResponse{Data: *issue})
}

// UpdateIssueHandler update issue handler
// @Summary create issue handler
// @Description FT update status and description for issue
// @Tags Issue
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param issue_id path string true "issue_id"
// @Param body body updateIssueRequest true "body"
// @Security ApiKeyAuth
// @Success 200 {object} issueResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /issues/{issue_id} [put]
func (h *Handler) UpdateIssueHandler(c *gin.Context) {
	request := updateIssueRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		util.HandleError(c, err)
		return
	}

	issueIdStr := c.Param("issue_id")
	issueId, err := strconv.Atoi(issueIdStr)
	if err != nil {
		h.log.Error("[handler.UpdateIssueHandler] strconv.Atoi(issueIdStr)", err)
		util.HandleError(c, errors.ErrInvalidIdFormat)
	}

	issue, err := h.doUpdateIssue(&request, issueId)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, issueResponse{Data: *issue})
}

// DeleteIssueByIdHandler delete issue by id handler
// @Summary delete issue handler
// @Description delete issue by issue id
// @Tags Issue
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param issue_id path string true "issue_id"
// @Security ApiKeyAuth
// @Success 200 {object} model.SuccessResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /issues/{issue_id} [delete]
func (h *Handler) DeleteIssueByIdHandler(c *gin.Context) {

	authInfo, err := util.GetAuthenInfoFromContext(c)
	if err != nil {
		h.log.Error("[handler.CreateIssueHandler] GetAuthenInfoFromContext()", err)
		util.HandleError(c, err)
		return
	}
	issueIdStr := c.Param("issue_id")
	issueId, err := strconv.Atoi(issueIdStr)
	if err != nil {
		h.log.Error("[handler.DeleteIssueByIdHandler] strconv.Atoi(issueIdStr)", err)
		util.HandleError(c, errors.ErrInvalidIdFormat)
	}

	err = h.doDeleteIssue(issueId, authInfo.UserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Code: http.StatusOK, Message: consts.SuccessMessage})
}

func (h *Handler) doCreateIssue(param *model.Issue) (*model.Issue, error) {
	_, err := h.repo.Site.GetById(h.store, param.SiteId)
	if err != nil {
		h.log.Error("[handler.doCreateIssue] Site.GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrSiteNotFound
		}
		return nil, errors.ErrInternalServerError
	}

	_, err = h.repo.Inspection.GetById(h.store, param.InspectionId)
	if err != nil {
		h.log.Error("[handler.doCreateIssue] Inspection.GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrInspectionNotFound
		}
		return nil, errors.ErrInternalServerError
	}

	tx, done := h.store.NewTransaction()

	issue, err := h.repo.Issue.Create(tx, *param)
	if err != nil {
		h.log.Error("[handler.doCreateIssue] Issue.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}

	var sparePartReq []model.IssueSparePart
	for _, v := range param.IssueSpareParts {
		sparePart, err := h.repo.SparePart.GetById(tx, v.SparePartId)
		if err != nil {
			h.log.Error("[handler.doCreateIssue] SparePart.GetById()", err)
			if inerr.Is(err, gorm.ErrRecordNotFound) {
				return nil, done(errors.ErrSparePartNotFound)
			}
			return nil, done(errors.ErrInternalServerError)
		}

		var iSparePart *model.IssueSparePart

		if sparePart.InStock >= v.Quantity {
			updatedQuantity := sparePart.InStock - v.Quantity
			_, err = h.repo.SparePart.UpdateQuantity(tx, int(sparePart.Id), updatedQuantity)
			if err != nil {
				h.log.Error("[handler.doCreateIssue] SparePart.UpdateQuantity()", err)
				return nil, done(errors.ErrInternalServerError)
			}
			iSparePart, err = h.repo.IssueSparePart.Create(tx, model.IssueSparePart{
				IssueId:     int(issue.Id),
				SparePartId: int(sparePart.Id),
				Quantity:    v.Quantity,
				Status:      consts.StatusCompleted})
			if err != nil {
				h.log.Error("[handler.doCreateIssue] IssueSparePart.Create()", err)
				return nil, done(errors.ErrInternalServerError)
			}
		} else {
			iSparePart, err = h.repo.IssueSparePart.Create(tx, model.IssueSparePart{
				IssueId:     int(issue.Id),
				SparePartId: int(sparePart.Id),
				Quantity:    v.Quantity,
				Status:      consts.StatusInProgress})
			if err != nil {
				h.log.Error("[handler.doCreateIssue] IssueSparePart.Create()", err)
				return nil, done(errors.ErrInternalServerError)
			}
		}
		sparePartReq = append(sparePartReq, *iSparePart)
	}
	issue.IssueSpareParts = sparePartReq
	newData, err := json.Marshal(issue)
	if err != nil {
		h.log.Error("[handler.doCreateIssue] json.Marshal(issue)", err)
		return nil, done(errors.ErrInternalServerError)
	}

	logAction := model.UserActionLog{
		UserID:     issue.UserId,
		TargetId:   int(issue.Id),
		TargetType: consts.TargetTypeIssue,
		Action:     consts.ActionDeleteIssue,
		NewData:    datatypes.JSON(newData),
		OldData:    datatypes.JSON(nil),
	}

	_, err = h.repo.UserActionLog.Create(tx, logAction)
	if err != nil {
		h.log.Error("[handler.doCreateIssue] UserActionLog.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}
	return issue, done(nil)
}

func (h *Handler) doUpdateIssue(param *updateIssueRequest, id int) (*model.Issue, error) {
	issue, err := h.repo.Issue.GetById(h.store, id)
	if err != nil {
		h.log.Error("[handler.doUpdateIssue] Issue.GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrIssueNotFound
		}
		return nil, errors.ErrInternalServerError
	}

	updatedIssueStruct := *issue
	updatedIssue := &updatedIssueStruct
	if param.Description != "" {
		updatedIssue.Description = param.Description
	}
	if param.Status == "completed" || param.Status == "in_progress" {
		updatedIssue.Status = param.Status
	}

	tx, done := h.store.NewTransaction()

	updatedIssue, err = h.repo.Issue.Update(tx, *updatedIssue)
	if err != nil {
		h.log.Error("[handler.doUpdateIssue] Issue.Update()", err)
		return nil, done(errors.ErrInternalServerError)
	}

	newData, err := json.Marshal(updatedIssue)
	if err != nil {
		h.log.Error("[handler.doUpdateIssue] json.Marshal(updatedIssue)", err)
		return nil, done(errors.ErrInternalServerError)
	}

	oldData, err := json.Marshal(issue)
	if err != nil {
		h.log.Error("[handler.doUpdateIssue] json.Marshal(issue)", err)
		return nil, done(errors.ErrInternalServerError)
	}

	logAction := model.UserActionLog{
		UserID:     issue.UserId,
		TargetId:   int(updatedIssue.Id),
		TargetType: consts.TargetTypeIssue,
		Action:     consts.ActionUpdateIssue,
		NewData:    datatypes.JSON(newData),
		OldData:    datatypes.JSON(oldData),
	}

	_, err = h.repo.UserActionLog.Create(tx, logAction)
	if err != nil {
		h.log.Error("[handler.doUpdateIssue] UserActionLog.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}
	return updatedIssue, done(nil)
}

func (h *Handler) doDeleteIssue(id int, userId int) error {
	issue, err := h.repo.Issue.GetById(h.store, id)
	if err != nil {
		h.log.Error("[handler.doDeleteIssue] Issue.GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return errors.ErrIssueNotFound
		}
		return errors.ErrInternalServerError
	}
	if issue.Status == consts.StatusCompleted {
		return errors.ErrIssueDeleteCompleted
	}

	tx, done := h.store.NewTransaction()

	for _, v := range issue.IssueSpareParts {
		if v.Status == consts.StatusCompleted {
			sparePart, err := h.repo.SparePart.GetById(tx, v.SparePartId)
			if err != nil {
				h.log.Error("[handler.doDeleteIssue] SparePart.GetById()", err)
				if inerr.Is(err, gorm.ErrRecordNotFound) {
					return done(errors.ErrSparePartNotFound)
				}
				return done(errors.ErrInternalServerError)
			}
			updatedQuantity := sparePart.InStock + v.Quantity
			_, err = h.repo.SparePart.UpdateQuantity(tx, int(sparePart.Id), updatedQuantity)
			if err != nil {
				h.log.Error("[handler.doDeleteIssue] SparePart.UpdateQuantity()", err)
				return done(errors.ErrInternalServerError)
			}
		}
	}

	err = h.repo.IssueSparePart.DeleteByIssueId(tx, int(issue.Id))
	if err != nil {
		h.log.Error("[handler.doDeleteIssue] IssueSparePart.DeleteByIssueId()", err)
		return done(errors.ErrInternalServerError)
	}

	err = h.repo.Issue.DeleteById(tx, int(issue.Id))
	if err != nil {
		h.log.Error("[handler.doDeleteIssue] Issue.DeleteById()", err)
		return done(errors.ErrInternalServerError)
	}

	oldData, err := json.Marshal(issue)
	if err != nil {
		h.log.Error("[handler.doDeleteIssue] json.Marshal(issue)", err)
		return done(errors.ErrInternalServerError)
	}

	logAction := model.UserActionLog{
		UserID:     userId,
		TargetId:   int(issue.Id),
		TargetType: consts.TargetTypeIssue,
		Action:     consts.ActionCreateIssue,
		NewData:    datatypes.JSON(nil),
		OldData:    datatypes.JSON(oldData),
	}

	_, err = h.repo.UserActionLog.Create(tx, logAction)
	if err != nil {
		h.log.Error("[handler.doDeleteIssue] UserActionLog.Create()", err)
		return done(errors.ErrInternalServerError)
	}
	return done(nil)
}
