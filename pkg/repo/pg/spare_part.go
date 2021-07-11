package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

type sparePartRepo struct{}

func NewSparePartRepo() repo.SparePartRepo {
	return &sparePartRepo{}
}

func (s *sparePartRepo) GetById(store repo.DBRepo, id int) (*model.SparePart, error) {
	db := store.DB()

	part := model.SparePart{}
	return &part, db.Where("id = ?", id).First(&part).Error
}

func (s *sparePartRepo) UpdateQuantity(store repo.DBRepo, id int, quantity int) (*model.SparePart, error) {
	db := store.DB()
	part := model.SparePart{}
	return &part, db.Model(&part).Where("id = ?", id).Update("in_stock", quantity).Error
}

func (s *sparePartRepo) GetAll(store repo.DBRepo) ([]model.SparePart, error) {
	db := store.DB()

	spareParts := []model.SparePart{}
	return spareParts, db.Find(&spareParts).Error
}
