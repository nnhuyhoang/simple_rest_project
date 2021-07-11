package mw

import (
	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/config"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/repo/pg"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
)

// AuthMiddleware ...
type AuthMiddleware struct {
	cfg   config.Config
	repo  *repo.Repo
	store repo.DBRepo
}

// AuthMiddleware ..
type AuthMiddleFunc interface {
	AhthPerm() gin.HandlerFunc
	WithAuth(c *gin.Context)
}

// NewAuthMiddleware ...
func NewAuthMiddleware(cfg config.Config, s repo.DBRepo) *AuthMiddleware {
	r := pg.NewRepo()
	return &AuthMiddleware{
		cfg:   cfg,
		store: s,
		repo:  r,
	}
}

func (a *AuthMiddleware) WithAuth(c *gin.Context) {
	token, err := authenticate(c, a.cfg)
	if err != nil {
		util.HandleError(c, err)
		c.Abort()
		return
	}

	authInfo, err := util.GetAuthInfoFromToken(token, a.cfg.JWTSecret)
	if err != nil {
		util.HandleError(c, err)
		c.Abort()
		return
	}
	c.Set(consts.AuthInfo, authInfo)
}

func authenticate(c *gin.Context, cfg config.Config) (string, error) {
	token, err := util.GetTokenFromContext(c)
	if err != nil {
		return "", err
	}

	return token, util.ValidateToken(token, cfg.JWTSecret)
}

func (a *AuthMiddleware) AuthPerm(roleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenInfo, err := util.GetAuthenInfoFromContext(c)
		if err != nil {
			util.HandleError(c, err)
			c.Abort()
			return
		}
		if authenInfo.RoleCode != roleCode {
			util.HandleError(c, errors.ErrPermissionDenied)
			c.Abort()
			return
		}
		c.Next()
	}
}
