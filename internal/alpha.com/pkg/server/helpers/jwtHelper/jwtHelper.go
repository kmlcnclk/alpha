package jwtHelper

import (
	"time"

	"alpha.com/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type IJwtHelper interface {
	CreateTokens(userID string) (string, string, error)
}

type jwtHelper struct {
}

func NewJwtHelper() IJwtHelper {
	return &jwtHelper{}
}

var jwtKey = []byte(configuration.JWT_SECRET)
var refreshKey = []byte(configuration.REFRESH_SECRET)

type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func (j *jwtHelper) CreateTokens(userID string) (string, string, error) {

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshExpirationTime := time.Now().Add(24 * time.Hour)
	refreshClaims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
			Id:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenString, nil
}
