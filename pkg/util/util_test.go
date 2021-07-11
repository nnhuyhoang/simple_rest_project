package util

import (
	"testing"
	"time"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
)

func TestValidateEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid email",
			args: args{
				email: "minh@gmail.com",
			},
			want: true,
		},
		{
			name: "invalid email",
			args: args{
				email: "minh.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.args.email); got != tt.want {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePasswordStringl(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "valid password",
			args: args{
				password: "password",
			},
			want: nil,
		},
		{
			name: "invalid password",
			args: args{
				password: "sdf",
			},
			want: errors.ErrInvalidPassword,
		},
		{
			name: "empty password",
			args: args{
				password: "",
			},
			want: errors.ErrInvalidPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePasswordString(tt.args.password); got != tt.want {
				t.Errorf("ValidatePasswordString() = %v, want nil", got)
			}
		})
	}
}

func TestParseDateRequest(t *testing.T) {
	timezone := consts.DefaultTimezone
	loc, _ := time.LoadLocation(timezone)
	validDate := time.Date(2021, time.Month(6), 13, 0, 0, 0, 0, loc)
	type args struct {
		date string
		tz   string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantValid bool
	}{
		{
			name: "valid date",
			args: args{
				date: "2021-06-13",
				tz:   timezone,
			},
			wantValid: true,
			wantErr:   false,
		},
		{
			name: "invalid date",
			args: args{
				date: "sdf-06-13",
				tz:   timezone,
			},
			wantValid: false,
			wantErr:   true,
		},
		{
			name: "invalid date",
			args: args{
				date: "06-13",
				tz:   timezone,
			},
			wantValid: false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateRequest(tt.args.date, tt.args.tz)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateRequest() = %v", err)
				return
			}
			if !tt.wantValid {
				return
			}
			if got.Before(validDate) || got.After(validDate) {
				t.Errorf("ParseDateRequest() = wrong date formatted")
				return
			}
		})
	}
}
