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

func TestHandler_doCreatePurchaseRequest(t *testing.T) {

	firstExistedSpart := 1
	secondExistedSpart := 2
	nonExistedSparePartId := 3

	userId := 1

	purchaseRequestId := 1

	firstISparePartId := 1
	secondISparePartId := 2
	nonExistedISpareParId := 3
	errExistedISpareParId := 4

	firstQuantity := 70
	secondQuantity := 60
	firstExtraQuantity := 50
	secondExtraQuantity := 40

	firstInStock := 100
	secondInStock := 150

	spareParts := []model.SparePart{
		{
			Id:      uint(firstExistedSpart),
			InStock: firstInStock,
			Name:    "name 1",
			Code:    "name_1",
		},
		{
			Id:      uint(secondExistedSpart),
			InStock: secondInStock,
			Name:    "name 2",
			Code:    "name_2",
		},
	}

	iSpareParts := []model.IssueSparePart{
		{
			Id:          uint(firstISparePartId),
			SparePartId: firstExistedSpart,
			Quantity:    firstQuantity,
			Status:      consts.StatusInProgress,
			SparePart:   &spareParts[0],
		},
		{
			Id:          uint(secondISparePartId),
			SparePartId: secondExistedSpart,
			Quantity:    secondQuantity,
			Status:      consts.StatusInProgress,
			SparePart:   &spareParts[1],
		},
	}

	iiSpareParts := []model.IssueSparePart{
		{
			Id:          uint(firstISparePartId),
			SparePartId: firstExistedSpart,
			Quantity:    firstQuantity,
			Status:      consts.StatusCompleted,
			SparePart:   &spareParts[0],
		},
		{
			Id:          uint(secondISparePartId),
			SparePartId: secondExistedSpart,
			Quantity:    secondQuantity,
			Status:      consts.StatusCompleted,
			SparePart:   &spareParts[1],
		},
	}

	eISparePart := model.IssueSparePart{
		Id:          uint(errExistedISpareParId),
		SparePartId: firstExistedSpart,
		Quantity:    firstQuantity,
		Status:      consts.StatusCompleted,
	}

	rPurchase_spare_parts := []model.PurchaseSparePart{
		{
			PurchaseRequestId: purchaseRequestId,
			SparePartId:       firstExistedSpart,
			Quantity:          firstQuantity + firstExtraQuantity,
			SparePartName:     "name 1",
			SparePartCode:     "name_1",
		},
		{
			PurchaseRequestId: purchaseRequestId,
			SparePartId:       secondExistedSpart,
			Quantity:          secondQuantity + secondExtraQuantity,
			SparePartName:     "name 2",
			SparePartCode:     "name_2",
		},
	}

	nPurchase_spare_parts := []model.PurchaseSparePart{
		{
			Id:                1,
			PurchaseRequestId: purchaseRequestId,
			SparePartId:       firstExistedSpart,
			Quantity:          firstQuantity + firstExtraQuantity,
		},
		{
			Id:                2,
			PurchaseRequestId: purchaseRequestId,
			SparePartId:       secondExistedSpart,
			Quantity:          secondQuantity + secondExtraQuantity,
		},
	}

	nPurchaseRequest := model.PurchaseRequest{
		Id:                 uint(purchaseRequestId),
		UserId:             userId,
		Status:             consts.StatusCompleted,
		PurchaseSpareParts: nPurchase_spare_parts,
	}

	userAction := model.UserActionLog{
		Id: 1,
	}

	firstIncreaseInStock := firstInStock + firstQuantity + firstExtraQuantity
	secondIncreaseInStock := secondInStock + secondQuantity + secondExtraQuantity

	firstDecreaseInStock := firstInStock + firstExtraQuantity
	secondDecreaseInStock := secondInStock + secondExtraQuantity

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userActionRepo := mock_repo.NewMockUserActionLogRepo(ctrl)
	issueSparePartRepo := mock_repo.NewMockIssueSparePartRepo(ctrl)
	sparePartRepo := mock_repo.NewMockSparePartRepo(ctrl)
	purchaseRequestRepo := mock_repo.NewMockPurchaseRequestRepo(ctrl)
	purchaseSparePartRepo := mock_repo.NewMockPurchaseSparePartRepo(ctrl)

	userActionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&userAction, nil).AnyTimes()

	purchaseRequestRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&nPurchaseRequest, nil).AnyTimes()

	issueSparePartRepo.EXPECT().GetbyId(gomock.Any(), firstISparePartId).Return(&iSpareParts[0], nil).AnyTimes()
	issueSparePartRepo.EXPECT().GetbyId(gomock.Any(), secondISparePartId).Return(&iSpareParts[1], nil).AnyTimes()
	issueSparePartRepo.EXPECT().GetbyId(gomock.Any(), errExistedISpareParId).Return(&eISparePart, nil).AnyTimes()
	issueSparePartRepo.EXPECT().GetbyId(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	issueSparePartRepo.EXPECT().Update(gomock.Any(), iiSpareParts[0]).Return(&model.IssueSparePart{Id: uint(firstISparePartId), Status: consts.StatusCompleted}, nil).AnyTimes()
	issueSparePartRepo.EXPECT().Update(gomock.Any(), iiSpareParts[1]).Return(&model.IssueSparePart{Id: uint(secondISparePartId), Status: consts.StatusCompleted}, nil).AnyTimes()

	sparePartRepo.EXPECT().GetById(gomock.Any(), firstExistedSpart).Return(&spareParts[0], nil).AnyTimes()
	sparePartRepo.EXPECT().GetById(gomock.Any(), secondExistedSpart).Return(&spareParts[1], nil).AnyTimes()
	sparePartRepo.EXPECT().GetById(gomock.Any(), nonExistedSparePartId).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	sparePartRepo.EXPECT().UpdateQuantity(gomock.Any(), firstExistedSpart, firstIncreaseInStock).Return(&model.SparePart{Id: uint(firstExistedSpart), InStock: firstIncreaseInStock}, nil).AnyTimes()
	sparePartRepo.EXPECT().UpdateQuantity(gomock.Any(), secondExistedSpart, secondIncreaseInStock).Return(&model.SparePart{Id: uint(secondExistedSpart), InStock: secondIncreaseInStock}, nil).AnyTimes()
	sparePartRepo.EXPECT().UpdateQuantity(gomock.Any(), firstExistedSpart, firstDecreaseInStock).Return(&model.SparePart{Id: uint(firstExistedSpart), InStock: firstDecreaseInStock}, nil).AnyTimes()
	sparePartRepo.EXPECT().UpdateQuantity(gomock.Any(), firstExistedSpart, secondDecreaseInStock).Return(&model.SparePart{Id: uint(secondExistedSpart), InStock: secondDecreaseInStock}, nil).AnyTimes()

	purchaseSparePartRepo.EXPECT().Create(gomock.Any(), rPurchase_spare_parts[0]).Return(&nPurchase_spare_parts[0], nil).AnyTimes()
	purchaseSparePartRepo.EXPECT().Create(gomock.Any(), rPurchase_spare_parts[1]).Return(&nPurchase_spare_parts[1], nil).AnyTimes()

	r := pg.NewRepo()
	r.IssueSparePart = issueSparePartRepo
	r.SparePart = sparePartRepo
	r.PurchaseSparePart = purchaseSparePartRepo
	r.PurchaseRequest = purchaseRequestRepo
	r.UserActionLog = userActionRepo
	h := NewTestHandler(r)

	type args struct {
		param *purchaseOrderRequest
		uId   int
	}
	tests := []struct {
		name    string
		args    args
		want    *model.PurchaseRequest
		wantErr bool
	}{
		{
			name: "no orders",
			args: args{
				param: &purchaseOrderRequest{Description: "something"},
				uId:   userId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "issue spare part not found",
			args: args{
				param: &purchaseOrderRequest{Description: "something", IssueOrderIds: []int{nonExistedISpareParId}},
				uId:   userId,
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "spare part not found",
			args: args{
				param: &purchaseOrderRequest{Description: "something", ExtraOrders: []extraOrder{{SparePartId: nonExistedSparePartId, Quantity: 30}}},
				uId:   userId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "issue spare part already completed",
			args: args{
				param: &purchaseOrderRequest{Description: "someting", IssueOrderIds: []int{errExistedISpareParId}},
				uId:   userId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			args: args{
				param: &purchaseOrderRequest{Description: "someting",
					IssueOrderIds: []int{firstExistedSpart, secondExistedSpart},
					ExtraOrders:   []extraOrder{{SparePartId: firstExistedSpart, Quantity: firstExtraQuantity}, {SparePartId: secondExistedSpart, Quantity: secondExtraQuantity}}},
				uId: userId,
			},
			want:    &nPurchaseRequest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.doCreatePurchaseRequest(tt.args.param, tt.args.uId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doCreatePurchaseRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.doCreateIssue() erorr = %v, want %v", got, tt.want)
				return
			}
		})
	}

}
