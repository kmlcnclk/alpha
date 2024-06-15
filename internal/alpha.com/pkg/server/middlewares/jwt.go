package middlewares

import (
	"context"
	"fmt"
	"strings"

	"alpha.com/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type UserContext struct {
	UserID string
}

var JWTSecret = []byte(configuration.JWT_SECRET)

func JwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing or malformed JWT")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header format")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired JWT")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userID"].(string)
		userCtx := &UserContext{UserID: userID}
		ctx := context.WithValue(c.UserContext(), "user", userCtx)
		c.SetUserContext(ctx)
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired JWT")
	}

	return c.Next()
}
