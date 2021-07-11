package pg

import (
	"time"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
	"gorm.io/gorm/clause"
)

type inspectionRepo struct{}

func NewInspectionRepo() repo.InspectionRepo {
	return &inspectionRepo{}
}

func (i *inspectionRepo) Create(s repo.DBRepo, param model.Inspection) (*model.Inspection, error) {
	return &param, s.DB().Select("SiteId", "UserId", "Date", "Status", "Remark", "CreatedAt", "UpdatedAt").Create(&param).Error
}

func (i *inspectionRepo) Update(s repo.DBRepo, param model.Inspection) (*model.Inspection, error) {
	return &param, s.DB().Omit(clause.Associations).Save(&param).Error
}

func (i *inspectionRepo) GetById(s repo.DBRepo, id int) (*model.Inspection, error) {
	db := s.DB()

	inspection := model.Inspection{}
	return &inspection, db.Preload("Site").Preload("Issues.IssueSpareParts").Where("id = ?", id).First(&inspection).Error
}

func (i *inspectionRepo) GetByUserIdWithDate(s repo.DBRepo, date time.Time, userId int) ([]model.Inspection, error) {
	tomorrow := date.AddDate(0, 0, 1)
	db := s.DB()
	inspections := []model.Inspection{}

	return inspections, db.Preload("Site").Preload("Issues", "status = ?", consts.StatusInProgress).Preload("Issues.IssueSpareParts").Where("user_id = ? AND Date >= ? AND Date < ?", userId, date, tomorrow).Find(&inspections).Error
}
