package handler

import (
	inerr "errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
	"gorm.io/gorm"
)

type getSiteRequest struct {
	Date string `json:"date" form:"date" example:"2021-06-15"`
}

type siteResponse struct {
	Data model.Site `json:"data"`
}

type siteListResponse struct {
	Data []model.Site `json:"data"`
}

// GetSiteByDateHandler get site by date handler
// @Summary get site by date handler
// @Description Manager dashboard call this
// @Tags Site
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param date query string true "2021-06-13"
// @Security ApiKeyAuth
// @Success 200 {object} siteListResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /sites [get]
func (h *Handler) GetSiteByDateHandler(c *gin.Context) {
	request := getSiteRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		util.HandleError(c, err)
		return
	}
	date, err := util.ParseDateRequest(request.Date, h.cfg.Timezone)
	if err != nil {
		h.log.Error("[handler.GetSiteByDateHandler] ParseDateRequest()", err)
		util.HandleError(c, err)
		return
	}

	sites, err := h.doGetSiteByDate(date)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, siteListResponse{Data: sites})
}

// GetSiteByIdHandler get site by id handler
// @Summary get site by id with date handler
// @Description site detail for a day
// @Tags Site
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param site_id path string true "site_id"
// @Param date query string true "2021-06-13"
// @Security ApiKeyAuth
// @Success 200 {object} siteListResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /sites/{site_id} [get]
func (h *Handler) GetSiteByIdHandler(c *gin.Context) {
	request := getSiteRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		util.HandleError(c, err)
		return
	}
	date, err := util.ParseDateRequest(request.Date, h.cfg.Timezone)
	if err != nil {
		h.log.Error("[handler.GetSiteByIdHandler] ParseDateRequest()", err)
		util.HandleError(c, err)
		return
	}
	siteIdStr := c.Param("site_id")
	siteId, err := strconv.Atoi(siteIdStr)
	if err != nil {
		h.log.Error("[handler.GetSiteByIdHandler] strconv.Atoi(siteIdStr)", err)
		util.HandleError(c, errors.ErrInvalidIdFormat)
		return
	}

	site, err := h.doGetSiteById(siteId, date)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, siteResponse{Data: *site})
}

func (h *Handler) doGetSiteByDate(date time.Time) ([]model.Site, error) {
	sites, err := h.repo.Site.GetByDate(h.store, date)
	if err != nil {
		h.log.Error("[handler.doGetSiteByDate] Site.GetByDate()", err)
		return nil, errors.ErrInternalServerError
	}

	var updatedSites []model.Site
	for _, site := range sites {
		status := consts.StatusCompleted
		for _, inspect := range site.Inspections {
			if inspect.Status == consts.StatusInProgress {
				status = consts.StatusInProgress
			}
		}
		if len(site.Issues) > 0 {
			status = consts.StatusMalfunctioned
		}
		site.Status = status
		updatedSites = append(updatedSites, site)
	}
	return updatedSites, nil
}

func (h *Handler) doGetSiteById(siteId int, date time.Time) (*model.Site, error) {
	site, err := h.repo.Site.GetByIdWithDate(h.store, siteId, date)
	if err != nil {
		h.log.Error("[handler.doGetSiteById] Site.GetByIdWithDate()", err)
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrSiteNotFound
		}

		return nil, errors.ErrInternalServerError
	}
	status := consts.StatusCompleted
	for _, inspect := range site.Inspections {
		if inspect.Status == consts.StatusInProgress {
			status = consts.StatusInProgress
		}
	}
	if len(site.Issues) > 0 {
		status = consts.StatusMalfunctioned
	}
	site.Status = status
	return site, nil
}
