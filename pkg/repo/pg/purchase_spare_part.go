package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

type purchaseSparePartRepo struct{}

func NewPurchaseSparePartRepo() repo.PurchaseSparePartRepo {
	return &purchaseSparePartRepo{}
}

func (i *purchaseSparePartRepo) Create(s repo.DBRepo, param model.PurchaseSparePart) (*model.PurchaseSparePart, error) {
	return &param, s.DB().Select("PurchaseRequestId", "SparePartId", "Quantity", "SparePartName", "SparePartCode", "CreatedAt", "UpdatedAt").Create(&param).Error
}
