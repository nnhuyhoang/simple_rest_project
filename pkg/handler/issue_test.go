package handler

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	mock_repo "github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/mocks"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/pg"
	"gorm.io/gorm"
)

func TestHandler_doCreateIssue(t *testing.T) {

	existedSiteId := 1
	nonExistedSiteId := 2

	existedInspectionId := 1
	nonExistedInspectionId := 2

	firstExistedSpart := 1
	secondExistedSpart := 2
	nonExistedSparePartId := 3

	userId := 1

	issueId := 1

	firstISparePartId := 1
	secondISparePartId := 2

	enoughQuantity := 70
	notEnoughQuantity := 60

	iDescription := "nothing"

	site := model.Site{
		Id: uint(existedSiteId),
	}

	inspection := model.Inspection{
		Id:     uint(existedInspectionId),
		SiteId: existedSiteId,
		UserId: userId,
	}

	rISparePart := []model.IssueSparePart{
		{
			IssueId:     issueId,
			SparePartId: firstExistedSpart,
			Quantity:    enoughQuantity,
			Status:      consts.StatusCompleted,
		},
		{
			IssueId:     issueId,
			SparePartId: secondExistedSpart,
			Quantity:    notEnoughQuantity,
			Status:      consts.StatusInProgress,
		},
	}

	nISparePart := []model.IssueSparePart{
		{
			Id:          uint(firstISparePartId),
			IssueId:     issueId,
			SparePartId: firstExistedSpart,
			Quantity:    enoughQuantity,
			Status:      consts.StatusCompleted,
		},
		{
			Id:          uint(secondISparePartId),
			IssueId:     issueId,
			SparePartId: secondExistedSpart,
			Quantity:    notEnoughQuantity,
			Status:      consts.StatusInProgress,
		},
	}

	rIssue := model.Issue{
		SiteId:          existedSiteId,
		InspectionId:    existedInspectionId,
		Status:          consts.StatusInProgress,
		Description:     iDescription,
		IssueSpareParts: rISparePart,
	}

	nIssue := model.Issue{
		Id:              uint(issueId),
		SiteId:          existedSiteId,
		InspectionId:    existedInspectionId,
		Status:          consts.StatusInProgress,
		Description:     iDescription,
		IssueSpareParts: nISparePart,
	}

	spareParts := []model.SparePart{
		{
			Id:      uint(firstExistedSpart),
			InStock: 100,
		},
		{
			Id:      uint(secondExistedSpart),
			InStock: 50,
		},
	}

	remainInstock := spareParts[0].InStock - enoughQuantity

	userAction := model.UserActionLog{
		Id: 1,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	siteRepo := mock_repo.NewMockSiteRepo(ctrl)
	inspectionRepo := mock_repo.NewMockInspectionRepo(ctrl)
	userActionRepo := mock_repo.NewMockUserActionLogRepo(ctrl)
	issueRepo := mock_repo.NewMockIssueRepo(ctrl)
	issueSparePartRepo := mock_repo.NewMockIssueSparePartRepo(ctrl)
	sparePartRepo := mock_repo.NewMockSparePartRepo(ctrl)

	siteRepo.EXPECT().GetById(gomock.Any(), existedSiteId).Return(&site, nil).AnyTimes()
	siteRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()

	inspectionRepo.EXPECT().GetById(gomock.Any(), existedInspectionId).Return(&inspection, nil).AnyTimes()
	inspectionRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()

	issueRepo.EXPECT().Create(gomock.Any(), rIssue).Return(&nIssue, nil).AnyTimes()

	issueSparePartRepo.EXPECT().Create(gomock.Any(), rISparePart[0]).Return(&nISparePart[0], nil).AnyTimes()
	issueSparePartRepo.EXPECT().Create(gomock.Any(), rISparePart[1]).Return(&nISparePart[1], nil).AnyTimes()

	sparePartRepo.EXPECT().GetById(gomock.Any(), firstExistedSpart).Return(&spareParts[0], nil).AnyTimes()
	sparePartRepo.EXPECT().GetById(gomock.Any(), secondExistedSpart).Return(&spareParts[1], nil).AnyTimes()
	sparePartRepo.EXPECT().UpdateQuantity(gomock.Any(), firstExistedSpart, remainInstock).Return(&model.SparePart{InStock: remainInstock}, nil).AnyTimes()

	userActionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&userAction, nil).AnyTimes()

	r := pg.NewRepo()
	r.Inspection = inspectionRepo
	r.Issue = issueRepo
	r.Site = siteRepo
	r.IssueSparePart = issueSparePartRepo
	r.SparePart = sparePartRepo
	r.UserActionLog = userActionRepo
	h := NewTestHandler(r)

	type args struct {
		issue *model.Issue
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Issue
		wantErr bool
	}{
		{
			name: "create issue with non-existed site",
			args: args{
				issue: &model.Issue{
					SiteId:          nonExistedSiteId,
					InspectionId:    existedInspectionId,
					Status:          consts.StatusInProgress,
					Description:     iDescription,
					IssueSpareParts: rISparePart,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "create issue with non-existed inspection",
			args: args{
				issue: &model.Issue{
					SiteId:          existedSiteId,
					InspectionId:    nonExistedInspectionId,
					Status:          consts.StatusInProgress,
					Description:     iDescription,
					IssueSpareParts: rISparePart,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "create issue with non-existed spare part",
			args: args{
				issue: &model.Issue{
					SiteId:       existedSiteId,
					InspectionId: nonExistedInspectionId,
					Status:       consts.StatusInProgress,
					Description:  iDescription,
					IssueSpareParts: []model.IssueSparePart{
						{
							IssueId:     issueId,
							SparePartId: nonExistedSparePartId,
							Quantity:    enoughQuantity,
						},
						{
							IssueId:     issueId,
							SparePartId: secondExistedSpart,
							Quantity:    notEnoughQuantity,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			args: args{
				issue: &rIssue,
			},
			want:    &nIssue,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.doCreateIssue(tt.args.issue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doCreateIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.doCreateIssue() erorr = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestHandler_doUpdateIssue(t *testing.T) {
	existedIssueId := 1
	nonExistedIssueId := 2
	siteId := 1
	userId := 2
	inspection_id := 3
	description := "nothing wrong"
	issue := model.Issue{
		Id:           uint(existedIssueId),
		SiteId:       siteId,
		UserId:       userId,
		InspectionId: inspection_id,
		Status:       "in_progress",
		Description:  "",
	}

	uIssue := model.Issue{
		Id:           uint(existedIssueId),
		SiteId:       siteId,
		UserId:       userId,
		InspectionId: inspection_id,
		Status:       consts.StatusCompleted,
		Description:  description,
	}

	userAction := model.UserActionLog{
		Id: 1,
	}

	updateRequest := updateIssueRequest{
		Status:      consts.StatusCompleted,
		Description: description,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	issueRepo := mock_repo.NewMockIssueRepo(ctrl)
	userActionRepo := mock_repo.NewMockUserActionLogRepo(ctrl)

	issueRepo.EXPECT().GetById(gomock.Any(), existedIssueId).Return(&issue, nil).AnyTimes()
	issueRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	issueRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&uIssue, nil).AnyTimes()
	userActionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&userAction, nil).AnyTimes()

	r := pg.NewRepo()
	r.Issue = issueRepo
	r.UserActionLog = userActionRepo
	h := NewTestHandler(r)

	type args struct {
		issue *updateIssueRequest
		iId   int
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Issue
		wantErr bool
	}{
		{
			name: "update issue with non-existed issue",
			args: args{
				issue: &updateRequest,
				iId:   nonExistedIssueId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			args: args{
				issue: &updateRequest,
				iId:   existedIssueId,
			},
			want:    &uIssue,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.doUpdateIssue(tt.args.issue, tt.args.iId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doUpdateIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.doUpdateIssue() erorr = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestHandler_doDeleteIssue(t *testing.T) {

	existedSiteId := 1

	existedInspectionId := 1

	firstExistedSpart := 1
	secondExistedSpart := 2
	nonExistedSparePartId := 3

	userId := 1

	issueId := 1
	nonExistedIssueId := 2
	errIssueId := 3

	firstISparePartId := 1
	secondISparePartId := 2

	enoughQuantity := 70
	notEnoughQuantity := 60

	iDescription := "nothing"

	ISparePart := []model.IssueSparePart{
		{
			Id:          uint(firstISparePartId),
			IssueId:     issueId,
			SparePartId: firstExistedSpart,
			Quantity:    enoughQuantity,
			Status:      consts.StatusCompleted,
		},
		{
			Id:          uint(secondISparePartId),
			IssueId:     issueId,
			SparePartId: secondExistedSpart,
			Quantity:    notEnoughQuantity,
			Status:      consts.StatusInProgress,
		},
	}
	eSparePart := []model.IssueSparePart{
		{
			Id:          uint(firstISparePartId),
			IssueId:     issueId,
			SparePartId: nonExistedSparePartId,
			Quantity:    enoughQuantity,
			Status:      consts.StatusCompleted,
		},
		{
			Id:          uint(secondISparePartId),
			IssueId:     issueId,
			SparePartId: secondExistedSpart,
			Quantity:    notEnoughQuantity,
			Status:      consts.StatusInProgress,
		},
	}

	issue := model.Issue{
		Id:              uint(issueId),
		SiteId:          existedSiteId,
		InspectionId:    existedInspectionId,
		Status:          consts.StatusInProgress,
		Description:     iDescription,
		IssueSpareParts: ISparePart,
	}

	errIssue := model.Issue{
		Id:              uint(errIssueId),
		SiteId:          existedSiteId,
		InspectionId:    existedInspectionId,
		Status:          consts.StatusInProgress,
		Description:     iDescription,
		IssueSpareParts: eSparePart,
	}

	spareParts := []model.SparePart{
		{
			Id:      uint(firstExistedSpart),
			InStock: 100,
		},
		{
			Id:      uint(secondExistedSpart),
			InStock: 50,
		},
	}

	updatedInStock := spareParts[0].InStock + enoughQuantity

	userAction := model.UserActionLog{
		Id: 1,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userActionRepo := mock_repo.NewMockUserActionLogRepo(ctrl)
	issueRepo := mock_repo.NewMockIssueRepo(ctrl)
	issueSparePartRepo := mock_repo.NewMockIssueSparePartRepo(ctrl)
	sparePartRepo := mock_repo.NewMockSparePartRepo(ctrl)

	issueRepo.EXPECT().GetById(gomock.Any(), issueId).Return(&issue, nil).AnyTimes()
	issueRepo.EXPECT().GetById(gomock.Any(), nonExistedIssueId).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	issueRepo.EXPECT().GetById(gomock.Any(), errIssueId).Return(&errIssue, nil).AnyTimes()
	issueRepo.EXPECT().DeleteById(gomock.Any(), issueId).Return(nil).AnyTimes()

	sparePartRepo.EXPECT().GetById(gomock.Any(), firstExistedSpart).Return(&spareParts[0], nil).AnyTimes()
	sparePartRepo.EXPECT().GetById(gomock.Any(), secondExistedSpart).Return(&spareParts[1], nil).AnyTimes()
	sparePartRepo.EXPECT().GetById(gomock.Any(), nonExistedSparePartId).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	sparePartRepo.EXPECT().UpdateQuantity(gomock.Any(), firstExistedSpart, updatedInStock).Return(&model.SparePart{InStock: updatedInStock}, nil)

	issueSparePartRepo.EXPECT().DeleteByIssueId(gomock.Any(), issueId).Return(nil).AnyTimes()

	userActionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&userAction, nil).AnyTimes()

	r := pg.NewRepo()
	r.Issue = issueRepo
	r.SparePart = sparePartRepo
	r.IssueSparePart = issueSparePartRepo
	r.UserActionLog = userActionRepo
	h := NewTestHandler(r)

	type args struct {
		iId int
		uId int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete issue with non-existed issue",
			args: args{
				iId: nonExistedIssueId,
				uId: userId,
			},
			wantErr: true,
		},
		{
			name: "create issue with non-existed spare parts",
			args: args{
				iId: errIssueId,
				uId: userId,
			},
			wantErr: true,
		},

		{
			name: "happy case",
			args: args{
				iId: issueId,
				uId: userId,
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := h.doDeleteIssue(tt.args.iId, tt.args.uId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doCreateIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
