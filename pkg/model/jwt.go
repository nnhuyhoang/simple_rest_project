package model

import "github.com/dgrijalva/jwt-go"

type AuthenticationInfo struct {
	jwt.StandardClaims
	UserID   int    `json:"id"`
	Email    string `json:"email"`
	RoleCode string `json:"roleCode"`
}
