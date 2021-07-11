package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/config"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/logger"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/pg"
	"github.com/nnhuyhoang/simple_rest_project/backend/services"
)

type Handler struct {
	log   logger.Log
	cfg   config.Config
	repo  *repo.Repo
	store repo.DBRepo
	svs   *services.Services
}

func NewHandler(cfg config.Config, l logger.Log, s repo.DBRepo, svs *services.Services) (*Handler, error) {
	r := pg.NewRepo()

	return &Handler{
		log:   l,
		cfg:   cfg,
		store: s,
		repo:  r,
		svs:   svs,
	}, nil
}

func NewTestHandler(r *repo.Repo) *Handler {
	h := &Handler{
		store: repo.NewTestStore(),
		log:   logger.NewJSONLogger(),
		repo:  pg.NewRepo(),
	}
	if r != nil {
		h.repo = r
	}
	return h
}

// Healthz handler
// Return "OK"
func (h *Handler) Healthz(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("OK"))
}
