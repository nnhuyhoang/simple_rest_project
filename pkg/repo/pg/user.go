package pg

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
)

type userRepo struct{}

func NewUserRepo() repo.UserRepo {
	return &userRepo{}
}

func (u *userRepo) Create(s repo.DBRepo, param model.User) (*model.User, error) {
	db := s.DB()

	return &param, db.Select("FirstName", "LastName", "FullName", "Email", "PhoneNumber", "HashedPassword", "RoleId", "CreatedAt", "UpdatedAt").Create(&param).Error
}

func (u *userRepo) GetAll(s repo.DBRepo) ([]model.User, error) {
	db := s.DB()

	users := []model.User{}
	return users, db.Preload("Role").Find(&users).Error
}

func (u *userRepo) GetByEmail(s repo.DBRepo, email string) (*model.User, error) {
	db := s.DB()

	user := model.User{}
	return &user, db.Preload("Role").Where("email = ?", email).First(&user).Error
}

func (u *userRepo) GetByPhone(s repo.DBRepo, phone string) (*model.User, error) {
	db := s.DB()

	user := model.User{}
	return &user, db.Preload("Role").Where("phone_number = ?", phone).First(&user).Error
}

func (u *userRepo) GetById(s repo.DBRepo, id int) (*model.User, error) {
	db := s.DB()

	user := model.User{}
	return &user, db.Where("id = ?", id).First(&user).Error
}

func (u *userRepo) GetByRoleCode(s repo.DBRepo, roleCode string) ([]model.User, error) {

	users := []model.User{}
	query := s.DB().
		Table("users").
		Joins("INNER JOIN roles ON users.role_id = roles.id").
		Where("roles.code = ?", roleCode).
		Preload("Role").
		Select("users.*")
	return users, query.Find(&users).Error
}
