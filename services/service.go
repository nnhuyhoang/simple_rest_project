package services

import (
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/config"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/logger"
	"github.com/nnhuyhoang/simple_rest_project/backend/services/email"
)

type Services struct {
	EmailService *email.EmailService
}

func NewServices(cfg config.Config, l logger.Log) *Services {
	return &Services{
		EmailService: email.NewEmailService(cfg, l),
	}
}
