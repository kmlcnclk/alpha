package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"alpha.com/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type IJwtService interface {
	CreateTokens(userID string) (string, string, error)
	ParseRefreshToken(refresh string) (string, error)
	CreateAccessToken(userID string) (string, error)
	CreateRefreshToken(userID string) (string, error)
}

type jwtService struct {
}

func NewJwtService() IJwtService {
	return &jwtService{}
}

var JWTSecret = []byte(configuration.JWT_SECRET)
var RefreshSecret = []byte(configuration.REFRESH_SECRET)

var AccessTokenTime = configuration.ACCESS_TOKEN_TIME
var RefreshTokenTime = configuration.REFRESH_TOKEN_TIME

type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func (j *jwtService) CreateTokens(userID string) (string, string, error) {
	accessToken, err := j.CreateAccessToken(userID)

	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := j.CreateRefreshToken(userID)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenString, nil
}

func (j *jwtService) ParseRefreshToken(refresh string) (string, error) {
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return RefreshSecret, nil
	})

	if err != nil {
		return "", errors.New("invalid or expired JWT")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userID"].(string)
		return userID, nil
	}

	return "", errors.New("invalid or expired JWT")
}

func (j *jwtService) CreateAccessToken(userID string) (string, error) {
	expirationDuration, err := parseDuration(AccessTokenTime)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(expirationDuration)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (j *jwtService) CreateRefreshToken(userID string) (string, error) {
	expirationDuration, err := parseDuration(RefreshTokenTime)
	if err != nil {
		return "", err
	}

	refreshExpirationTime := time.Now().Add(expirationDuration)
	refreshClaims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
			Id:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(RefreshSecret)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func parseDuration(durationStr string) (time.Duration, error) {
	var duration time.Duration

	numStr := durationStr[:len(durationStr)-1]
	unitStr := strings.ToLower(string(durationStr[len(durationStr)-1]))

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return duration, fmt.Errorf("invalid duration format: %v", err)
	}

	switch unitStr {
	case "s":
		duration = time.Duration(num) * time.Second
	case "m":
		duration = time.Duration(num) * time.Minute
	case "h":
		duration = time.Duration(num) * time.Hour
	case "d":
		duration = time.Duration(num) * 24 * time.Hour
	default:
		return duration, fmt.Errorf("unknown duration unit: %v", unitStr)
	}

	return duration, nil
}
