package pg

import "github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"

func NewRepo() *repo.Repo {
	return &repo.Repo{
		User:              NewUserRepo(),
		Role:              NewRoleRepo(),
		Inspection:        NewInspectionRepo(),
		Site:              NewSiteRepo(),
		SparePart:         NewSparePartRepo(),
		Issue:             NewIssueRepo(),
		IssueSparePart:    NewIssueSparePartRepo(),
		UserActionLog:     NewUserActionLogRepo(),
		PurchaseRequest:   NewPurchaseRequestRepo(),
		PurchaseSparePart: NewPurchaseSparePartRepo(),
	}
}
