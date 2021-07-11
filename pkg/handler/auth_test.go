package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	mock_repo "github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/mocks"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/pg"
)

func TestHandler_doLoginWithMail(t *testing.T) {
	validPassword := "password"
	hashedPassword := "$2y$12$fJ9IyDoJ70KABsjHOdTSuu4hO9F0mfXLBx2uC6Oszb98MdMfoRJt2 "
	invalidPassword := "RandomInvalidPassword"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := pg.NewRepo()
	u := mock_repo.NewMockUserRepo(ctrl)

	r.User = u
	h := NewTestHandler(r)

	u.EXPECT().GetByEmail(gomock.Any(), "jon@gmail.com").Return(&model.User{
		Email:          "jon@gmail.com",
		HashedPassword: hashedPassword,
		Role: &model.Role{
			Name: "Field Technician",
			Code: "field_technician",
		},
	}, nil).AnyTimes()

	tests := []struct {
		name    string
		input   loginRequest
		wantErr bool
	}{
		{
			name: "Login with mail success",
			input: loginRequest{
				Email:    "jon@gmail.com",
				Password: validPassword,
			},
			wantErr: false,
		},
		{
			name: "Login with mail failed, wrong password",
			input: loginRequest{
				Email:    "jon@gmail.com",
				Password: invalidPassword,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := h.doLoginWithMail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.doLoginWithMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (response.Data.AccessToken == "") {
				t.Error("Fail to login by mail")
				return
			}
		})
	}
}
