package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/nnhuyhoang/simple_rest_project/backend/docs/swagger"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/config"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/handler"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/logger"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/routes"
	"github.com/nnhuyhoang/simple_rest_project/backend/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Usecase2B API
// @version 1.0
// @description This is a Usecase2B server.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cls := config.DefaultConfigLoaders()
	cfg := config.LoadConfig(cls)

	l := initLog(cfg)

	s := repo.NewPostgresStore(&cfg)

	svs := services.NewServices(cfg, l)
	defer svs.EmailService.Channel.Close()

	router := setupRouter(cfg, l, s, svs)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

}

func initLog(cfg config.Config) logger.Log {
	return logger.NewJSONLogger(
		logger.WithServiceName(cfg.ServiceName),
		logger.WithHostName(cfg.BaseURL),
	)
}

func setupRouter(cfg config.Config, l logger.Log, s repo.DBRepo, svs *services.Services) *gin.Engine {
	r := gin.Default()

	h, err := handler.NewHandler(cfg, l, s, svs)
	if err != nil {
		log.Fatal(err)
	}

	r.Use(cors.New(
		cors.Config{
			AllowOrigins: cfg.GetCORS(),
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
			AllowHeaders: []string{"Origin", "Host",
				"Content-Type", "Content-Length",
				"Accept-Encoding", "Accept-Language", "Accept",
				"X-CSRF-Token", "Authorization", "X-Requested-With", "X-Access-Token"},
			ExposeHeaders:    []string{"MeAllowMethodsntent-Length"},
			AllowCredentials: true,
		},
	))
	url := ginSwagger.URL(cfg.BaseURL + "/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// handlers
	r.GET("/healthz", h.Healthz)
	routes.NewRoutes(r, h, cfg, s)
	return r
}
