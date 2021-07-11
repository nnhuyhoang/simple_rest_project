package handler

import (
	"encoding/json"
	inerr "errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type createInspectionRequest struct {
	SiteId int    `json:"siteId" form:"siteId" binding:"required" example:"1"`
	UserId int    `json:"userId" form:"userId" binding:"required" example:"2"`
	Date   string `json:"date" form:"date" binding:"required" example:"2021-06-15"`
}

type updateInspectionRequest struct {
	Status string `json:"status" form:"status" binding:"required" example:"completed"`
	Remark string `json:"remark" form:"remark" example:"nothing wrong"`
}

type getInspectionByUserIdRequest struct {
	Date string `json:"date" form:"date" example:"2021-06-15"`
}

type inspectionResponse struct {
	Data model.Inspection `json:"data"`
}

type inspectionListResponse struct {
	Data []model.Inspection `json:"data"`
}

// CreateInspectionHandler create inspection handler
// @Summary create inspection handler
// @Description Manager will inspection for FT
// @Tags Inspection
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param body body handler.createInspectionRequest true "create inspect form"
// @Security ApiKeyAuth
// @Success 200 {object} inspectionResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /inspections [post]
func (h *Handler) CreateInspectionHandler(c *gin.Context) {

	authInfo, err := util.GetAuthenInfoFromContext(c)
	if err != nil {
		h.log.Error("[handler.CreateInspectionHandler] GetAuthenInfoFromContext()", err)
		util.HandleError(c, err)
		return
	}
	request := createInspectionRequest{}
	if err = c.ShouldBindJSON(&request); err != nil {
		util.HandleError(c, err)
		return
	}

	inspectDate, err := util.ParseDateRequest(request.Date, h.cfg.Timezone)
	if err != nil {
		h.log.Error("[handler.CreateInspectionHandler] ParseDateRequest()", err)
		util.HandleError(c, err)
		return
	}
	inspectionParam := model.Inspection{
		SiteId: request.SiteId,
		UserId: request.UserId,
		Date:   inspectDate,
		Status: consts.StatusInProgress,
	}

	inspection, err := h.doCreateInspection(&inspectionParam, authInfo.UserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, inspectionResponse{Data: *inspection})
}

// UpdateInspectionHandler update inspection handler
// @Summary update inspection handler
// @Description FT put status and remark after inspection
// @Tags Inspection
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param inspection_id path string true "inspection_id"
// @Param body body handler.updateInspectionRequest true "update inspect form"
// @Security ApiKeyAuth
// @Success 200 {object} inspectionResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /inspections/{inspection_id} [put]
func (h *Handler) UpdateInspectionHandler(c *gin.Context) {
	authInfo, err := util.GetAuthenInfoFromContext(c)
	if err != nil {
		h.log.Error("[handler.UpdateInspectionHandler] GetAuthenInfoFromContext()", err)
		util.HandleError(c, err)
		return
	}
	request := updateInspectionRequest{}
	if err = c.ShouldBindJSON(&request); err != nil {
		util.HandleError(c, err)
		return
	}
	inspectionIdStr := c.Param("inspection_id")
	inspectionId, err := strconv.Atoi(inspectionIdStr)
	if err != nil {
		h.log.Error("[handler.UpdateInspectionHandler] strconv.Atoi(inspectionIdStr)", err)
		util.HandleError(c, errors.ErrInvalidIdFormat)
	}
	inspection, err := h.doUpdateInspection(&request, inspectionId, authInfo.UserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, inspectionResponse{Data: *inspection})
}

// GetInspectionByUserIdHandler get inspection by userId handler
// @Summary get inspection by userId handler
// @Description FT get list of sites need to inspect by userId and date
// @Tags Inspection
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param date query string true "2021-06-13"
// @Security ApiKeyAuth
// @Success 200 {object} inspectionListResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /inspections [get]
func (h *Handler) GetInspectionByUserIdHandler(c *gin.Context) {
	authInfo, err := util.GetAuthenInfoFromContext(c)
	if err != nil {
		h.log.Error("[handler.GetInspectionByUserId] GetAuthenInfoFromContext()", err)
		util.HandleError(c, err)
		return
	}

	request := getInspectionByUserIdRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		util.HandleError(c, err)
		return
	}
	inspectDate, err := util.ParseDateRequest(request.Date, h.cfg.Timezone)
	if err != nil {
		h.log.Error("[handler.GetInspectionByUserId] ParseDateRequest()", err)
		util.HandleError(c, err)
		return
	}

	inspections, err := h.repo.Inspection.GetByUserIdWithDate(h.store, inspectDate, authInfo.UserID)
	if err != nil {
		h.log.Error("[handler.GetInspectionByUserId] GetByUserIdWithDate()", err)
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, inspectionListResponse{Data: inspections})
}

// GetInspectionByIdHandler get inspection by Id handler
// @Summary get inspection by Id handler
// @Description get inspection by id
// @Tags Inspection
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param inspection_id path string true "inspection_id"
// @Security ApiKeyAuth
// @Success 200 {object} inspectionListResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /inspections/{inspection_id} [get]
func (h *Handler) GetInspectionByIdHandler(c *gin.Context) {
	inspectionIdStr := c.Param("inspection_id")
	inspectionId, err := strconv.Atoi(inspectionIdStr)
	if err != nil {
		h.log.Error("[handler.GetInspectionByIdHandler] strconv.Atoi(inspectionIdStr)", err)
		util.HandleError(c, errors.ErrInvalidIdFormat)
		return
	}

	inspection, err := h.repo.Inspection.GetById(h.store, inspectionId)
	if err != nil {
		h.log.Error("[handler.GetInspectionByUserId] GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			util.HandleError(c, errors.ErrInspectionNotFound)
			return
		}
		util.HandleError(c, errors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, inspectionResponse{Data: *inspection})
}

func (h *Handler) doCreateInspection(param *model.Inspection, created_by int) (*model.Inspection, error) {
	_, err := h.repo.User.GetById(h.store, param.UserId)
	if err != nil {
		h.log.Error("[handler.doCreateInspection] User.GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServerError
	}

	_, err = h.repo.Site.GetById(h.store, param.SiteId)
	if err != nil {
		h.log.Error("[handler.doCreateInspection] Site.GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrSiteNotFound
		}
		return nil, errors.ErrInternalServerError
	}

	tx, done := h.store.NewTransaction()

	inspection, err := h.repo.Inspection.Create(tx, *param)
	if err != nil {
		h.log.Error("[handler.doCreateInspection] Inspection.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}
	newData, err := json.Marshal(inspection)
	if err != nil {
		h.log.Error("[handler.doUpdateInspection] json.Marshal(inspection)", err)
		return nil, done(errors.ErrInternalServerError)
	}

	logAction := model.UserActionLog{
		UserID:     created_by,
		TargetId:   int(inspection.Id),
		TargetType: consts.TargetTypeInspection,
		Action:     consts.ActionCreateInspection,
		NewData:    datatypes.JSON(newData),
		OldData:    datatypes.JSON(nil),
	}

	_, err = h.repo.UserActionLog.Create(tx, logAction)
	if err != nil {
		h.log.Error("[handler.doCreateInspection] UserActionLog.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}
	return inspection, done(nil)
}

func (h *Handler) doUpdateInspection(param *updateInspectionRequest, inspection_id int, created_by int) (*model.Inspection, error) {
	inspection, err := h.repo.Inspection.GetById(h.store, inspection_id)
	if err != nil {
		h.log.Error("[handler.doUpdateInspection] Inspection.GetById()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrInspectionNotFound
		}
		return nil, errors.ErrInternalServerError
	}

	updatedInspectionStruct := *inspection

	updatedInspection := &updatedInspectionStruct
	updatedInspection.Status = param.Status
	updatedInspection.Remark = param.Remark
	updatedInspection.UpdatedAt = time.Now()

	tx, done := h.store.NewTransaction()
	updatedInspection, err = h.repo.Inspection.Update(tx, *updatedInspection)
	if err != nil {
		h.log.Error("[handler.doUpdateInspection] Inspection.Update()", err)
		return nil, done(errors.ErrInternalServerError)
	}

	newData, err := json.Marshal(updatedInspection)
	if err != nil {
		h.log.Error("[handler.doUpdateInspection] json.Marshal(updatedInspection)", err)
		return nil, done(errors.ErrInternalServerError)
	}

	oldData, err := json.Marshal(inspection)
	if err != nil {
		h.log.Error("[handler.doUpdateInspection] json.Marshal(inspection)", err)
		return nil, done(errors.ErrInternalServerError)
	}

	logAction := model.UserActionLog{
		UserID:     created_by,
		TargetId:   int(inspection.Id),
		TargetType: consts.TargetTypeInspection,
		Action:     consts.ActionUpdateInspection,
		NewData:    datatypes.JSON(newData),
		OldData:    datatypes.JSON(oldData),
	}
	_, err = h.repo.UserActionLog.Create(tx, logAction)
	if err != nil {
		h.log.Error("[handler.doUpdateInspection] UserActionLog.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}
	return updatedInspection, done(nil)
}
