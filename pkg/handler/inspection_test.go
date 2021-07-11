package handler

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	mock_repo "github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/mocks"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/pg"
	"gorm.io/gorm"
)

func TestHandler_doCreateInspection(t *testing.T) {
	loc, _ := time.LoadLocation(consts.DefaultTimezone)
	createdbyId := 1
	userId := 3
	nonExistedUserId := 4
	existedSiteId := 1
	nonExistedSiteId := 2
	date := time.Date(2021, time.Month(6), 13, 0, 0, 0, 0, loc)
	site := model.Site{
		Id:          uint(existedSiteId),
		Name:        "RandomName",
		Description: "Nothing",
	}

	user := model.User{
		Id: uint(userId),
	}

	rInspection := model.Inspection{
		SiteId: existedSiteId,
		UserId: userId,
		Date:   date,
		Status: consts.StatusInProgress,
		Remark: "",
	}

	rUserInspection := model.Inspection{
		SiteId: nonExistedSiteId,
		UserId: nonExistedUserId,
		Date:   date,
		Status: consts.StatusInProgress,
		Remark: "",
	}
	rSiteInspection := model.Inspection{
		SiteId: nonExistedSiteId,
		UserId: userId,
		Date:   date,
		Status: consts.StatusInProgress,
		Remark: "",
	}
	nInspection := model.Inspection{
		Id:     uint(userId),
		SiteId: existedSiteId,
		UserId: userId,
		Date:   date,
		Status: consts.StatusInProgress,
		Remark: "",
	}

	userAction := model.UserActionLog{
		Id: 1,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	siteRepo := mock_repo.NewMockSiteRepo(ctrl)
	inspectionRepo := mock_repo.NewMockInspectionRepo(ctrl)
	userRepo := mock_repo.NewMockUserRepo(ctrl)
	userActionRepo := mock_repo.NewMockUserActionLogRepo(ctrl)

	siteRepo.EXPECT().GetById(gomock.Any(), existedSiteId).Return(&site, nil).AnyTimes()
	siteRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	inspectionRepo.EXPECT().Create(gomock.Any(), rInspection).Return(&nInspection, nil).AnyTimes()
	userRepo.EXPECT().GetById(gomock.Any(), userId).Return(&user, nil).AnyTimes()
	userRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	userActionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&userAction, nil).AnyTimes()

	r := pg.NewRepo()
	r.Inspection = inspectionRepo
	r.User = userRepo
	r.Site = siteRepo
	r.UserActionLog = userActionRepo
	h := NewTestHandler(r)

	type args struct {
		inpspection *model.Inspection
		uID         int
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Inspection
		wantErr bool
	}{
		{
			name: "create inspection with non-existed user",
			args: args{
				inpspection: &rUserInspection,
				uID:         createdbyId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "create inspection with non-existed site",
			args: args{
				inpspection: &rSiteInspection,
				uID:         createdbyId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			args: args{
				inpspection: &rInspection,
				uID:         createdbyId,
			},
			want:    &nInspection,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.doCreateInspection(tt.args.inpspection, tt.args.uID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doCreateInspection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.doCreateInspection() erorr = %v, want %v", got, tt.want)
				return
			}
		})
	}

}

func TestHandler_doUpdateInspection(t *testing.T) {
	existedInspectionId := 1
	nonExistedInspectionId := 2
	siteId := 1
	userId := 2
	created_by := 3
	remark := "nothing wrong"
	inspection := model.Inspection{
		Id:     uint(existedInspectionId),
		SiteId: siteId,
		UserId: userId,
		Status: "in_progress",
		Remark: "",
	}

	uInspection := model.Inspection{
		Id:     uint(existedInspectionId),
		SiteId: siteId,
		UserId: userId,
		Status: consts.StatusCompleted,
		Remark: remark,
	}

	userAction := model.UserActionLog{
		Id: 1,
	}

	updateRequest := updateInspectionRequest{
		Status: consts.StatusCompleted,
		Remark: remark,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inspectionRepo := mock_repo.NewMockInspectionRepo(ctrl)
	userActionRepo := mock_repo.NewMockUserActionLogRepo(ctrl)

	inspectionRepo.EXPECT().GetById(gomock.Any(), existedInspectionId).Return(&inspection, nil).AnyTimes()
	inspectionRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	inspectionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&uInspection, nil).AnyTimes()
	userActionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&userAction, nil).AnyTimes()

	r := pg.NewRepo()
	r.Inspection = inspectionRepo
	r.UserActionLog = userActionRepo
	h := NewTestHandler(r)

	type args struct {
		inpspection *updateInspectionRequest
		uId         int
		iId         int
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Inspection
		wantErr bool
	}{
		{
			name: "update inspection with non-existed inpsect",
			args: args{
				inpspection: &updateRequest,
				uId:         created_by,
				iId:         nonExistedInspectionId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			args: args{
				inpspection: &updateRequest,
				uId:         created_by,
				iId:         existedInspectionId,
			},
			want:    &uInspection,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.doUpdateInspection(tt.args.inpspection, tt.args.iId, tt.args.uId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doUpdateInspection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.doUpdateInspection() erorr = %v, want %v", got, tt.want)
				return
			}
		})
	}

}
