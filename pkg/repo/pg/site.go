package pg

import (
	"time"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

type siteRepo struct{}

func NewSiteRepo() repo.SiteRepo {
	return &siteRepo{}
}

func (s *siteRepo) GetById(store repo.DBRepo, id int) (*model.Site, error) {
	db := store.DB()

	site := model.Site{}
	return &site, db.Where("id = ?", id).First(&site).Error
}

func (s *siteRepo) GetByIdWithDate(store repo.DBRepo, id int, date time.Time) (*model.Site, error) {
	tomorrow := date.AddDate(0, 0, 1)
	db := store.DB()
	site := model.Site{}
	query := db.Table("sites").
		Where("id = ?", id).
		Preload("Inspections", "date >= ? AND date < ?", date, tomorrow).
		Preload("Inspections.User").
		Preload("Issues", "status = ?", consts.StatusInProgress).
		Preload("Issues.IssueSpareParts")
	return &site, query.First(&site).Error
}

func (s *siteRepo) GetByDate(store repo.DBRepo, date time.Time) ([]model.Site, error) {
	tomorrow := date.AddDate(0, 0, 1)
	db := store.DB()
	sites := []model.Site{}
	query := db.Table("sites").
		Joins("INNER JOIN inspections ON inspections.site_id = sites.id").
		Where("inspections.date >= ? AND inspections.date < ?", date, tomorrow).
		Preload("Inspections", "date >= ? AND date < ?", date, tomorrow).
		Preload("Inspections.User").
		Preload("Issues", "status = ?", consts.StatusInProgress).
		Distinct("sites.*").
		Select("sites.*")
	return sites, query.Find(&sites).Error
}
