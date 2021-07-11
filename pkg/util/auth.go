package util

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

// ValidatePassword compare hashedPassword with password input is twin or not
func ValidatePassword(hashedPassword string, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GenerateJWTToken ...
func GenerateJWTToken(info *model.AuthenticationInfo, expiresAt int64, secretKey string) (string, error) {
	info.ExpiresAt = expiresAt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, info)
	encryptedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return encryptedToken, nil
}

// GetTokenFromContext ...
func GetTokenFromContext(c *gin.Context) (string, error) {
	headers := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(headers) != 2 {
		return "", errors.ErrUnexpectedHeader
	}
	if headers[0] != "Bearer" {
		return "", errors.ErrInvalidAuthenType
	}
	return headers[1], nil
}

// ValidateToken
func ValidateToken(accessToken string, secretKey string) error {
	claims := &model.AuthenticationInfo{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return errors.ErrInvalidAccessToken
	}

	return claims.Valid()
}

func GetAuthInfoFromToken(accessToken string, secretKey string) (*model.AuthenticationInfo, error) {
	claims := model.AuthenticationInfo{}
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if !token.Valid {
		return nil, errors.ErrInvalidAccessToken
	}
	if err != nil {
		return nil, errors.ErrInvalidAccessToken
	}
	return &claims, nil
}

// GetAuthenInfoFromContext ...
func GetAuthenInfoFromContext(c *gin.Context) (*model.AuthenticationInfo, error) {
	a, ok := c.Get(consts.AuthInfo)
	if !ok {
		return nil, errors.ErrInvalidAccessToken
	}

	authenInfo, ok := a.(*model.AuthenticationInfo)
	if !ok {
		return nil, errors.ErrInvalidAccessToken
	}

	return authenInfo, nil
}
