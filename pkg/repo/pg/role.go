package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

type roleRepo struct{}

func NewRoleRepo() repo.RoleRepo {
	return &roleRepo{}
}

func (r *roleRepo) GetByRoleCode(s repo.DBRepo, roleCode string) (*model.Role, error) {

	role := model.Role{}

	return &role, s.DB().Where("code = ?", roleCode).First(&role).Error
}
