package repo

import (
	"time"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
)

type Repo struct {
	User              UserRepo
	Role              RoleRepo
	Inspection        InspectionRepo
	Site              SiteRepo
	SparePart         SparePartRepo
	Issue             IssueRepo
	IssueSparePart    IssueSparePartRepo
	UserActionLog     UserActionLogRepo
	PurchaseRequest   PurchaseRequestRepo
	PurchaseSparePart PurchaseSparePartRepo
}

type UserRepo interface {
	Create(s DBRepo, param model.User) (*model.User, error)
	GetAll(s DBRepo) ([]model.User, error)
	GetByEmail(s DBRepo, email string) (*model.User, error)
	GetByPhone(s DBRepo, phone string) (*model.User, error)
	GetByRoleCode(s DBRepo, roleCode string) ([]model.User, error)
	GetById(s DBRepo, id int) (*model.User, error)
}

type InspectionRepo interface {
	GetById(s DBRepo, id int) (*model.Inspection, error)
	Create(s DBRepo, param model.Inspection) (*model.Inspection, error)
	GetByUserIdWithDate(s DBRepo, date time.Time, userId int) ([]model.Inspection, error)
	Update(s DBRepo, param model.Inspection) (*model.Inspection, error)
}

type SiteRepo interface {
	GetById(s DBRepo, id int) (*model.Site, error)
	GetByDate(s DBRepo, date time.Time) ([]model.Site, error)
	GetByIdWithDate(s DBRepo, id int, date time.Time) (*model.Site, error)
}

type UserActionLogRepo interface {
	Create(s DBRepo, param model.UserActionLog) (*model.UserActionLog, error)
}

type SparePartRepo interface {
	GetAll(s DBRepo) ([]model.SparePart, error)
	GetById(s DBRepo, id int) (*model.SparePart, error)
	UpdateQuantity(s DBRepo, id int, quantity int) (*model.SparePart, error)
}

type IssueRepo interface {
	Create(s DBRepo, param model.Issue) (*model.Issue, error)
	Update(s DBRepo, param model.Issue) (*model.Issue, error)
	DeleteById(s DBRepo, id int) error
	GetById(s DBRepo, id int) (*model.Issue, error)
	GetByIdWithSparePart(s DBRepo, id int) (*model.Issue, error)
}

type IssueSparePartRepo interface {
	GetByStatus(s DBRepo, status string) ([]model.IssueSparePart, error)
	Create(s DBRepo, param model.IssueSparePart) (*model.IssueSparePart, error)
	Update(s DBRepo, param model.IssueSparePart) (*model.IssueSparePart, error)
	GetbyId(s DBRepo, id int) (*model.IssueSparePart, error)
	DeleteByIssueId(s DBRepo, issueId int) error
}

type PurchaseRequestRepo interface {
	Create(s DBRepo, param model.PurchaseRequest) (*model.PurchaseRequest, error)
	GetAll(s DBRepo) ([]model.PurchaseRequest, error)
}

type PurchaseSparePartRepo interface {
	Create(s DBRepo, param model.PurchaseSparePart) (*model.PurchaseSparePart, error)
}

type RoleRepo interface {
	GetByRoleCode(s DBRepo, roleCode string) (*model.Role, error)
}
