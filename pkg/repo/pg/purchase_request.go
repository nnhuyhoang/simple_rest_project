package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

type purchaseRequestRepo struct{}

func NewPurchaseRequestRepo() repo.PurchaseRequestRepo {
	return &purchaseRequestRepo{}
}

func (i *purchaseRequestRepo) GetAll(s repo.DBRepo) ([]model.PurchaseRequest, error) {
	db := s.DB()

	puchaseRequests := []model.PurchaseRequest{}
	return puchaseRequests, db.Preload("PurchaseSpareParts").Order("order_date desc").Find(&puchaseRequests).Error
}

func (i *purchaseRequestRepo) Create(s repo.DBRepo, param model.PurchaseRequest) (*model.PurchaseRequest, error) {
	return &param, s.DB().Select("Description", "UserId", "OrderDate", "Status", "CreatedAt", "UpdatedAt").Create(&param).Error
}
