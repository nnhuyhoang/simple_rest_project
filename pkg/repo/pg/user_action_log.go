package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

type userActionLogRepo struct{}

func NewUserActionLogRepo() repo.UserActionLogRepo {
	return &userActionLogRepo{}
}

func (i *userActionLogRepo) Create(s repo.DBRepo, param model.UserActionLog) (*model.UserActionLog, error) {
	return &param, s.DB().Create(&param).Error
}
