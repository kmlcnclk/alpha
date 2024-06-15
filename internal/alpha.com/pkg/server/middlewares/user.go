package middlewares

import (
	"regexp"

	"alpha.com/internal/alpha.com/application/controller/request"
	"github.com/gofiber/fiber/v2"
)

func IsEmailFormatCorrect(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodPost {
		return ctx.Next()
	}

	var requestBody request.UserCreteRequest
	if err := ctx.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "JSON format is not correct")
	}

	if isEmailContent(requestBody.Email) {
		return ctx.Next()
	}

	return fiber.NewError(fiber.StatusBadRequest, "Email format is not correct")
}

func isEmailContent(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}
