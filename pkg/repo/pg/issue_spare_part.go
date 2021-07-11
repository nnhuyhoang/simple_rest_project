package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
	"gorm.io/gorm/clause"
)

type issueSparePartRepo struct{}

func NewIssueSparePartRepo() repo.IssueSparePartRepo {
	return &issueSparePartRepo{}
}

func (i *issueSparePartRepo) Create(s repo.DBRepo, param model.IssueSparePart) (*model.IssueSparePart, error) {
	return &param, s.DB().Select("IssueId", "SparePartId", "Quantity", "Status", "CreatedAt", "UpdatedAt").Create(&param).Error
}

func (i *issueSparePartRepo) DeleteByIssueId(s repo.DBRepo, issueId int) error {
	db := s.DB()

	part := model.IssueSparePart{}
	return db.Where("issue_id = ?", issueId).Delete(&part).Error
}

func (i *issueSparePartRepo) GetByStatus(s repo.DBRepo, status string) ([]model.IssueSparePart, error) {
	db := s.DB()
	parts := []model.IssueSparePart{}
	return parts, db.Where("status = ?", status).Find(&parts).Error
}

func (i *issueSparePartRepo) GetbyId(s repo.DBRepo, id int) (*model.IssueSparePart, error) {
	db := s.DB()
	part := model.IssueSparePart{}
	return &part, db.Preload("SparePart").Where("id = ?", id).First(&part).Error
}

func (i *issueSparePartRepo) Update(s repo.DBRepo, param model.IssueSparePart) (*model.IssueSparePart, error) {
	return &param, s.DB().Omit(clause.Associations).Save(&param).Error
}
