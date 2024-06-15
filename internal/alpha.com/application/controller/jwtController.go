package controller

import (
	"fmt"
	"net/http"

	"alpha.com/internal/alpha.com/application/controller/request"
	"alpha.com/internal/alpha.com/application/controller/response"
	"alpha.com/internal/alpha.com/application/handler/jwt"
	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type IJwtController interface {
	Create(ctx *fiber.Ctx) error
	GetJwt(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
}

type JwtController struct {
	jwtQueryService   query.IJwtQueryService
	jwtCommandHandler jwt.ICommandHandler
	customValidator   validation.ICustomValidator
}

func NewJwtController(jwtQueryService query.IJwtQueryService, jwtCommandHandler jwt.ICommandHandler, customValidator validation.ICustomValidator) IJwtController {
	return &JwtController{
		jwtQueryService:   jwtQueryService,
		jwtCommandHandler: jwtCommandHandler,
		customValidator:   customValidator,
	}
}

// Save godoc

//	@Summary		This method used for saving new jwt
//	@Description	saving new jwt
//
// @Param requestBody body request.JwtCreteRequest nil "Handle Request Body"
//
//	@Tags			JWT
//	@Accept			json
//	@Produce		json
//
// @Success 200
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/jwt [post]
func (u *JwtController) Create(ctx *fiber.Ctx) error {
	var req request.JwtCreteRequest
	err := ctx.BodyParser(&req)

	fmt.Printf("jwtController.Save INFO -> Request Body: %v\n", req)

	if err != nil {
		fmt.Printf("jwtController.Save ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	accessToken, refreshToken, err := u.jwtCommandHandler.Create(ctx.Context(), req.ToCommand())

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	)
}

// GetJwt godoc
//
//	@Summary		This method used for get all jwts
//	@Description	get all jwts
//	@Tags			JWT
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} []response.JwtResponse
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/jwt [get]
func (u *JwtController) GetJwt(ctx *fiber.Ctx) error {
	jwts, err := u.jwtQueryService.Get(ctx.Context())

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		response.ToJwtResponseList(jwts),
	)
}

// Save godoc

//	@Summary		This method used for saving new jwt
//	@Description	saving new jwt
//
// @Param x-refresh header string true "{token}"
//
//	@Tags			JWT
//	@Accept			json
//	@Produce		json
//
// @Success 200
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/jwt/refresh [post]
func (u *JwtController) Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Get("X-Refresh")

	if refreshToken == "" {
		return fiber.NewError(http.StatusBadRequest, "Refresh token is required")
	}

	fmt.Printf("jwtController.Refresh INFO -> Refresh Token: %v\n", refreshToken)

	userID, err := u.jwtQueryService.ParseRefreshToken(ctx.Context(), refreshToken)

	if err != nil {
		fmt.Printf("jwtController.Refresh ERROR -> There was an error while parsing refresh token - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	user, err := u.jwtQueryService.GetUserById(ctx.Context(), userID)

	if err != nil {
		fmt.Printf("jwtController.Refresh ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return fiber.NewError(http.StatusNotFound, "User not found with given refresh token")
	}

	accessToken, err := u.jwtCommandHandler.Refresh(ctx.Context(), userID)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = u.jwtCommandHandler.Update(ctx.Context(), userID, accessToken, refreshToken)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		map[string]interface{}{
			"accessToken": accessToken,
		},
	)
}
