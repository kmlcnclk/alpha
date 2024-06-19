package controller

import (
	"fmt"
	"net/http"

	"alpha.com/internal/alpha.com/application/controller/request"
	"alpha.com/internal/alpha.com/application/controller/response"
	"alpha.com/internal/alpha.com/application/handler/businessAccount"
	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/pkg/server/middlewares"
	"alpha.com/internal/alpha.com/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type IBusinessAccountController interface {
	Save(ctx *fiber.Ctx) error
	GetAllBusinessAccounts(ctx *fiber.Ctx) error
}

type BusinessAccountController struct {
	businessAccountQueryService   query.IBusinessAccountQueryService
	businessAccountCommandHandler businessAccount.ICommandHandler
	customValidator               validation.ICustomValidator
}

func NewBusinessAccountController(
	businessAccountQueryService query.IBusinessAccountQueryService,
	businessAccountCommandHandler businessAccount.ICommandHandler,
	customValidator validation.ICustomValidator,
) IBusinessAccountController {
	return &BusinessAccountController{
		businessAccountQueryService:   businessAccountQueryService,
		businessAccountCommandHandler: businessAccountCommandHandler,
		customValidator:               customValidator,
	}
}

// Save godoc

//	@Summary		This method used for saving new business account
//	@Description	saving new business account
//
//	@Tags			Business Accounts
//	@Accept			json
//	@Produce		json
//
// @Param requestBody body request.BusinessAccountCreateRequest nil "Handle Request Body"
//
// @Param Authorization header string true "Bearer {token}"
//
// @Success 200
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/business-account [post]
func (u *BusinessAccountController) Save(ctx *fiber.Ctx) error {
	var req request.BusinessAccountCreateRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		fmt.Printf("BusinessAccountController.Save ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("BusinessAccountController.Save STARTED with request: %#v\n", req)

	if err := u.customValidator.Validate(req); err != nil {
		fmt.Printf("BusinessAccountController.Save INVALID request: %#v - ERROR: %#v\n", req, err)
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	userCtx := ctx.UserContext().Value("user").(*middlewares.UserContext)

	errOfCommandHandler := u.businessAccountCommandHandler.Save(ctx.UserContext(), req.ToCommand(), userCtx.UserID)

	if errOfCommandHandler != nil {
		return fiber.NewError(http.StatusBadRequest, errOfCommandHandler.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		map[string]interface{}{
			"message": "Business Account Successfully Created",
		},
	)
}

// GetUser godoc
//
//	@Summary		This method used for get all business accounts
//	@Description	get all business accounts
//	@Tags			Business Accounts
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} []response.BusinessAccountResponse
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/business-account [get]
func (u *BusinessAccountController) GetAllBusinessAccounts(ctx *fiber.Ctx) error {
	businessAccounts, err := u.businessAccountQueryService.GetAllBusinessAccounts(ctx.UserContext())

	if err != nil {
		fmt.Printf("businessAccountController.GetAllBusinessAccounts ERROR -> There was an error while getting businessAccounts - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.ToBusinessAccountResponseList(businessAccounts))
}
