package handler

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	mock_repo "github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/mocks"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/pg"
	"gorm.io/gorm"
)

func TestHandler_doGetSiteByDate(t *testing.T) {

	iSiteId := 1
	cSiteId := 2
	mSiteId := 3

	var emptyIssues []model.Issue
	iInspection := []model.Inspection{
		{
			Id:     1,
			SiteId: iSiteId,
			Status: consts.StatusInProgress,
		},
		{
			Id:     2,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
	}
	iIssue := []model.Issue{
		{
			Id:     1,
			SiteId: iSiteId,
			Status: consts.StatusInProgress,
		},
		{
			Id:     2,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
	}
	cInspection := []model.Inspection{
		{
			Id:     1,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
		{
			Id:     2,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
	}
	iSite := model.Site{
		Id:          uint(iSiteId),
		Inspections: iInspection,
		Issues:      emptyIssues,
	}

	cSite := model.Site{
		Id:          uint(cSiteId),
		Inspections: cInspection,
		Issues:      emptyIssues,
	}

	mSite := model.Site{
		Id:          uint(mSiteId),
		Inspections: iInspection,
		Issues:      iIssue,
	}

	sites := []model.Site{
		iSite, cSite, mSite,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	siteRepo := mock_repo.NewMockSiteRepo(ctrl)

	siteRepo.EXPECT().GetByDate(gomock.Any(), gomock.Any()).Return(sites, nil).AnyTimes()

	r := pg.NewRepo()
	r.Site = siteRepo
	h := NewTestHandler(r)

	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want []model.Site
	}{
		{
			name: "happy case",
			args: args{
				date: time.Now(),
			},
			want: []model.Site{
				{
					Id:     uint(iSiteId),
					Status: consts.StatusInProgress,
				},
				{
					Id:     uint(cSiteId),
					Status: consts.StatusCompleted,
				},
				{
					Id:     uint(mSiteId),
					Status: consts.StatusMalfunctioned,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sites, _ := h.doGetSiteByDate(tt.args.date)
			for i, v := range sites {
				if v.Status != tt.want[i].Status {
					t.Errorf("Handler.doGetSiteByDate() error = %v, want = %v", v.Status, tt.want[i].Status)
					return
				}
			}

		})
	}
}

func TestHandler_doGetSiteById(t *testing.T) {

	iSiteId := 1
	cSiteId := 2
	mSiteId := 3
	errSiteId := 4

	var emptyIssues []model.Issue
	iInspection := []model.Inspection{
		{
			Id:     1,
			SiteId: iSiteId,
			Status: consts.StatusInProgress,
		},
		{
			Id:     2,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
	}
	iIssue := []model.Issue{
		{
			Id:     1,
			SiteId: iSiteId,
			Status: consts.StatusInProgress,
		},
		{
			Id:     2,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
	}
	cInspection := []model.Inspection{
		{
			Id:     1,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
		{
			Id:     2,
			SiteId: iSiteId,
			Status: consts.StatusCompleted,
		},
	}
	iSite := model.Site{
		Id:          uint(iSiteId),
		Inspections: iInspection,
		Issues:      emptyIssues,
	}

	cSite := model.Site{
		Id:          uint(cSiteId),
		Inspections: cInspection,
		Issues:      emptyIssues,
	}

	mSite := model.Site{
		Id:          uint(mSiteId),
		Inspections: iInspection,
		Issues:      iIssue,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	siteRepo := mock_repo.NewMockSiteRepo(ctrl)

	siteRepo.EXPECT().GetByIdWithDate(gomock.Any(), iSiteId, gomock.Any()).Return(&iSite, nil).AnyTimes()
	siteRepo.EXPECT().GetByIdWithDate(gomock.Any(), cSiteId, gomock.Any()).Return(&cSite, nil).AnyTimes()
	siteRepo.EXPECT().GetByIdWithDate(gomock.Any(), mSiteId, gomock.Any()).Return(&mSite, nil).AnyTimes()
	siteRepo.EXPECT().GetByIdWithDate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	r := pg.NewRepo()
	r.Site = siteRepo
	h := NewTestHandler(r)

	type args struct {
		sId  int
		date time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Site
		wantErr bool
	}{
		{
			name: "get in_progress site",
			args: args{
				date: time.Now(),
				sId:  iSiteId,
			},
			want: &model.Site{
				Id:     uint(iSiteId),
				Status: consts.StatusInProgress,
			},
			wantErr: false,
		},
		{
			name: "get completed site",
			args: args{
				date: time.Now(),
				sId:  cSiteId,
			},
			want: &model.Site{
				Id:     uint(cSiteId),
				Status: consts.StatusCompleted,
			},
			wantErr: false,
		},
		{
			name: "get malfuntioned site",
			args: args{
				date: time.Now(),
				sId:  mSiteId,
			},
			want: &model.Site{
				Id:     uint(mSiteId),
				Status: consts.StatusMalfunctioned,
			},
			wantErr: false,
		},
		{
			name: "get non-existed site",
			args: args{
				date: time.Now(),
				sId:  errSiteId,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site, err := h.doGetSiteById(tt.args.sId, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doGetSiteById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil {
				return
			}
			if site.Status != tt.want.Status {
				t.Errorf("Handler.doGetSiteByDate() error = %v, want = %v", site.Status, tt.want.Status)
				return
			}

		})
	}
}
