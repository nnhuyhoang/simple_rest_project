package repo

import (
	"fmt"
	"net/http"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/config"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBRepo interface {
	DB() *gorm.DB
	NewTransaction() (DBRepo, FinallyFunc)
}

type store struct {
	Database *gorm.DB
}

type FinallyFunc = func(error) error

func (s *store) DB() *gorm.DB {
	return s.Database
}

func (s *store) NewTransaction() (newRepo DBRepo, finallyFunc FinallyFunc) {
	newDB := s.Database.Begin()

	finallyFunc = func(err error) error {
		if err != nil {
			nErr := newDB.Rollback().Error
			if nErr != nil {
				return errors.NewStringError(nErr.Error(), http.StatusInternalServerError)
			}
			return err
		}
		cErr := newDB.Commit().Error
		if cErr != nil {
			return errors.NewStringError(cErr.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	return &store{Database: newDB}, finallyFunc
}

func NewPostgresStore(cfg *config.Config) DBRepo {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	return &store{Database: db}
}
