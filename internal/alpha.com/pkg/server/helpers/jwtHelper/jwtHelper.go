package jwtHelper

import (
	"errors"
	"fmt"
	"time"

	"alpha.com/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type IJwtHelper interface {
	CreateTokens(userID string) (string, string, error)
	ParseRefreshToken(refresh string) (string, error)
	CreateAccessToken(userID string) (string, error)
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

func (j *jwtHelper) ParseRefreshToken(refresh string) (string, error) {
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshKey, nil
	})

	if err != nil {
		return "", errors.New("Invalid or expired JWT")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userID"].(string)
		return userID, nil
	}

	return "", errors.New("Invalid or expired JWT")
}

func (j *jwtHelper) CreateAccessToken(userID string) (string, error) {
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
		return "", err
	}

	return accessToken, nil
}
