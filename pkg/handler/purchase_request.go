package handler

import (
	"encoding/json"
	inerr "errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type purchaseOrderRequest struct {
	Description   string       `json:"description" form:"description" example:"order for something"`
	IssueOrderIds []int        `json:"issueOrderIds" form:"issueOrderIds"`
	ExtraOrders   []extraOrder `json:"extraOrders" form:"extraOrders"`
}

type extraOrder struct {
	SparePartId int `json:"sparePartId" form:"sparePartId" binding:"required" `
	Quantity    int `json:"quantity" form:"quantity" binding:"required"`
}

type order struct {
	SparePartId   int    `json:"sparePartId"`
	Quantity      int    `json:"quantity"`
	InStock       int    `json:"inStock"`
	SparePartName string `json:"sparePartName"`
	SparePartCode string `json:"sparePartCode"`
}

type purchaseOrderResponse struct {
	Data model.PurchaseRequest `json:"purchaseRequest"`
}
type purchaseOrderListResponse struct {
	Data []model.PurchaseRequest `json:"purchaseRequest"`
}

// GetAllPurchaseRequestHandler get all purchase request handler
// @Summary get all purchase request handler
// @Description get all purchase requests
// @Tags PurchaseRequest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Security ApiKeyAuth
// @Success 200 {object} purchaseOrderListResponse	"ok"
// @Failure 500 {object} model.Error
// @Router /purchase-requests [get]
func (h *Handler) GetAllPurchaseRequestHandler(c *gin.Context) {
	purchaseRequests, err := h.repo.PurchaseRequest.GetAll(h.store)
	if err != nil {
		h.log.Error("[handler.GetAllPurchaseRequestHandler] PurchaseRequest.GetAll()", err)
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, purchaseOrderListResponse{Data: purchaseRequests})
}

// CreatePurchaseRequestHandler create purchase request handler
// @Summary create puchase request handler
// @Description Manager will create puchase request
// @Tags PurchaseRequest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer "
// @Param body body handler.purchaseOrderRequest true "create purchase request form"
// @Security ApiKeyAuth
// @Success 200 {object} purchaseOrderResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /purchase-requests [post]
func (h *Handler) CreatePurchaseRequestHandler(c *gin.Context) {
	authInfo, err := util.GetAuthenInfoFromContext(c)
	if err != nil {
		h.log.Error("[handler.CreatePurchaseRequestHandler] GetAuthenInfoFromContext()", err)
		util.HandleError(c, err)
		return
	}

	request := purchaseOrderRequest{}
	if err = c.ShouldBindJSON(&request); err != nil {
		util.HandleError(c, err)
		return
	}

	purchaseRequest, err := h.doCreatePurchaseRequest(&request, authInfo.UserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, purchaseOrderResponse{Data: *purchaseRequest})
}

func (h *Handler) doCreatePurchaseRequest(param *purchaseOrderRequest, userId int) (*model.PurchaseRequest, error) {
	if len(param.IssueOrderIds) <= 0 && len(param.ExtraOrders) <= 0 {
		h.log.Error("[handler.doCreatePurchaseRequest] len(param.IssueOrderIds)() || en(param.ExtraOrders)")
		return nil, errors.ErrNoOrderFound
	}

	tx, done := h.store.NewTransaction()

	purchaseRequest, err := h.repo.PurchaseRequest.Create(tx, model.PurchaseRequest{
		Description: param.Description,
		UserId:      userId,
		OrderDate:   time.Now(),
		Status:      consts.StatusCompleted,
	})
	if err != nil {
		h.log.Error("[handler.doCreatePurchaseRequest] PurchaseRequest.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}

	var SparePartOrders []order
	var IssueSpareParts []model.IssueSparePart
	for _, orderId := range param.IssueOrderIds {
		iPart, err := h.repo.IssueSparePart.GetbyId(tx, orderId)
		if err != nil {
			h.log.Error("[handler.doCreatePurchaseRequest] IssueSparePart.GetbyId()", err)
			if inerr.Is(err, gorm.ErrRecordNotFound) {
				return nil, done(errors.ErrIssueSparePartNotFound)
			}
			return nil, done(errors.ErrInternalServerError)
		}
		if iPart.Status == consts.StatusCompleted {
			return nil, done(errors.ErrOrderAlreadyCompleted)
		}

		IssueSpareParts = append(IssueSpareParts, *iPart)

		index := -1
		for i, v := range SparePartOrders {
			if v.SparePartId == iPart.SparePartId {
				index = i
				break
			}
		}
		if index >= 0 {
			SparePartOrders[index].Quantity = SparePartOrders[index].Quantity + iPart.Quantity
		} else {
			SparePartOrders = append(SparePartOrders, order{
				SparePartId:   iPart.SparePartId,
				Quantity:      iPart.Quantity,
				InStock:       iPart.SparePart.InStock,
				SparePartName: iPart.SparePart.Name,
				SparePartCode: iPart.SparePart.Code})
		}
	}
	for _, extra := range param.ExtraOrders {
		iSPart, err := h.repo.SparePart.GetById(tx, extra.SparePartId)
		if err != nil {
			h.log.Error("[handler.doCreatePurchaseRequest] SparePart.GetbyId()", err)
			if inerr.Is(err, gorm.ErrRecordNotFound) {
				return nil, done(errors.ErrSparePartNotFound)
			}
			return nil, done(errors.ErrInternalServerError)
		}
		index := -1
		for i, v := range SparePartOrders {
			if v.SparePartId == int(iSPart.Id) {
				index = i
				break
			}
		}
		if index >= 0 {
			SparePartOrders[index].Quantity = SparePartOrders[index].Quantity + extra.Quantity
		} else {
			SparePartOrders = append(SparePartOrders, order{
				SparePartId:   int(iSPart.Id),
				Quantity:      extra.Quantity,
				InStock:       iSPart.InStock,
				SparePartName: iSPart.Name,
				SparePartCode: iSPart.Code})
		}
	}

	var createdOrders []model.PurchaseSparePart
	var updatedSpareParts []model.SparePart
	for _, order := range SparePartOrders {
		createdOrder, err := h.repo.PurchaseSparePart.Create(tx, model.PurchaseSparePart{
			PurchaseRequestId: int(purchaseRequest.Id),
			Quantity:          order.Quantity,
			SparePartId:       order.SparePartId,
			SparePartName:     order.SparePartName,
			SparePartCode:     order.SparePartCode,
		})
		if err != nil {
			h.log.Error("[handler.doCreatePurchaseRequest] PurchaseSparePart.Create()", err)
			return nil, done(errors.ErrInternalServerError)
		}
		createdOrders = append(createdOrders, *createdOrder)
		updatedInStock := order.InStock + createdOrder.Quantity
		updatedSparePart, err := h.repo.SparePart.UpdateQuantity(tx, createdOrder.SparePartId, updatedInStock)
		if err != nil {
			h.log.Error("[handler.doCreatePurchaseRequest] SparePart.UpdateQuantity()", err)
			return nil, done(errors.ErrInternalServerError)
		}
		updatedSparePart.Id = uint(createdOrder.SparePartId)
		updatedSpareParts = append(updatedSpareParts, *updatedSparePart)
	}

	purchaseRequest.PurchaseSpareParts = createdOrders

	for _, part := range IssueSpareParts {
		part.Status = consts.StatusCompleted
		updatedPart, err := h.repo.IssueSparePart.Update(tx, part)
		if err != nil {
			h.log.Error("[handler.doCreatePurchaseRequest] IssueSparePart.Update()", err)
			return nil, done(errors.ErrInternalServerError)
		}
		for i, spare := range updatedSpareParts {
			if updatedPart.SparePartId == int(spare.Id) {
				updatedQuantity := spare.InStock - updatedPart.Quantity
				_, err = h.repo.SparePart.UpdateQuantity(tx, int(spare.Id), updatedQuantity)
				if err != nil {
					h.log.Error("[handler.doCreatePurchaseRequest] SparePart.UpdateQuantity()-2", err)
					return nil, done(errors.ErrInternalServerError)
				}
				updatedSpareParts[i].InStock = updatedQuantity
				break
			}
		}
	}

	newData, err := json.Marshal(purchaseRequest)
	if err != nil {
		h.log.Error("[handler.doCreatePurchaseRequest] json.Marshal(purchaseRequest)", err)
		return nil, done(errors.ErrInternalServerError)
	}

	logAction := model.UserActionLog{
		UserID:     userId,
		TargetId:   int(purchaseRequest.Id),
		TargetType: consts.TargetPurchaseRequest,
		Action:     consts.ActionCreatePurchaseRequest,
		NewData:    datatypes.JSON(newData),
		OldData:    datatypes.JSON(nil),
	}

	_, err = h.repo.UserActionLog.Create(tx, logAction)
	if err != nil {
		h.log.Error("[handler.doCreatePurchaseRequest] UserActionLog.Create()", err)
		return nil, done(errors.ErrInternalServerError)
	}
	return purchaseRequest, done(nil)
}
