package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
	"gorm.io/gorm/clause"
)

type issueRepo struct{}

func NewIssueRepo() repo.IssueRepo {
	return &issueRepo{}
}

func (i *issueRepo) Create(s repo.DBRepo, param model.Issue) (*model.Issue, error) {
	return &param, s.DB().Select("SiteId", "UserId", "InspectionId", "Status", "Description", "CreatedAt", "UpdatedAt").Create(&param).Error
}

func (i *issueRepo) Update(s repo.DBRepo, param model.Issue) (*model.Issue, error) {
	return &param, s.DB().Omit(clause.Associations).Save(&param).Error
}

func (i *issueRepo) DeleteById(s repo.DBRepo, id int) error {
	db := s.DB()

	issue := model.Issue{}
	return db.Where("Id = ?", id).Delete(&issue).Error
}

func (s *issueRepo) GetById(store repo.DBRepo, id int) (*model.Issue, error) {
	db := store.DB()

	issue := model.Issue{}
	return &issue, db.Preload("IssueSpareParts").Where("id = ?", id).First(&issue).Error
}

func (s *issueRepo) GetByIdWithSparePart(store repo.DBRepo, id int) (*model.Issue, error) {
	db := store.DB()

	issue := model.Issue{}
	return &issue, db.Preload("IssueSpareParts.SparePart").Where("id = ?", id).First(&issue).Error
}
