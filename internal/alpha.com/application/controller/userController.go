package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"alpha.com/configuration"
	"alpha.com/internal/alpha.com/application/controller/request"
	"alpha.com/internal/alpha.com/application/controller/response"
	"alpha.com/internal/alpha.com/application/handler/user"
	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/pkg/server/helpers"
	"alpha.com/internal/alpha.com/pkg/server/middlewares"
	"alpha.com/internal/alpha.com/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type IUserController interface {
	Save(ctx *fiber.Ctx) error
	GetUser(ctx *fiber.Ctx) error
	GetUserById(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
}

type UserController struct {
	userQueryService   query.IUserQueryService
	userCommandHandler user.ICommandHandler
	customValidator    validation.ICustomValidator
}

func NewUserController(userQueryService query.IUserQueryService, userCommandHandler user.ICommandHandler, customValidator validation.ICustomValidator) IUserController {
	return &UserController{
		userQueryService:   userQueryService,
		userCommandHandler: userCommandHandler,
		customValidator:    customValidator,
	}
}

// Save godoc

//	@Summary		This method used for saving new user
//	@Description	saving new user
//
// @Param requestBody body request.UserCreateRequest nil "Handle Request Body"
//
//	@Tags			User
//	@Accept			json
//	@Produce		json
//
// @Success 200
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/user [post]
func (u *UserController) Save(ctx *fiber.Ctx) error {
	var req request.UserCreateRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		fmt.Printf("userController.Save ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("userController.Save STARTED with request: %#v\n", req)

	if err := u.customValidator.Validate(req); err != nil {
		fmt.Printf("userController.Save INVALID request: %#v - ERROR: %#v\n", req, err)
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	userID, errOfCommandHandler := u.userCommandHandler.Save(ctx.UserContext(), req.ToCommand())

	if errOfCommandHandler != nil {
		return fiber.NewError(http.StatusBadRequest, errOfCommandHandler.Error())
	}

	if userID == "" {
		fmt.Printf("userController.Save ERROR -> There was an error while binding json - ERROR: %v\n", "Internal Server Error")
		return fiber.NewError(http.StatusBadRequest, "Internal Server Error")
	}

	// Create a map to represent the JSON data
	requestData := map[string]string{"userId": userID}

	// Convert the map to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	body, errOfHelper := helpers.HttpPostHelper(configuration.BACKEND_URL+"/api/v1/alpha/jwt", bytes.NewBuffer(jsonData))

	if errOfHelper != nil {
		fmt.Printf("userController.Save ERROR -> There was an error while sending request to jwt - ERROR: %v\n", errOfHelper.Error())
		return fiber.NewError(http.StatusInternalServerError, errOfHelper.Error())
	}

	if body == nil {
		fmt.Printf("userController.Save ERROR -> There was an error while sending request to jwt - ERROR: %v\n", "body is nil")
		return fiber.NewError(http.StatusBadRequest, "Body is nil")
	}

	var data map[string]string

	errOfData := json.Unmarshal([]byte(string(body)), &data)
	if errOfData != nil {
		fmt.Printf("userController.Save ERROR -> There was an error while binding data - ERROR: %v\n:", errOfData.Error())
		return fiber.NewError(http.StatusInternalServerError, errOfData.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		map[string]interface{}{
			"message": "User Created Successfully",
			"response": map[string]interface{}{
				"accessToken":  data["accessToken"],
				"refreshToken": data["refreshToken"],
			},
		},
	)
}

// Save godoc

//	@Summary		This method used for sign up to user
//	@Description	sign up for user
//
// @Param requestBody body request.UserSignInRequest nil "Handle Request Body"
//
//	@Tags			User
//	@Accept			json
//	@Produce		json
//
// @Success 200
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/user/sign-in [post]
func (u *UserController) SignIn(ctx *fiber.Ctx) error {
	var req request.UserSignInRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		fmt.Printf("userController.SignIn ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("userController.SignIn STARTED with request: %#v\n", req)

	if err := u.customValidator.Validate(req); err != nil {
		fmt.Printf("userController.SignIn INVALID request: %#v - ERROR: %#v\n", req, err)
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	userID, errOfCommandHandler := u.userCommandHandler.SignIn(ctx.UserContext(), req.ToCommand())

	if errOfCommandHandler != nil {
		return fiber.NewError(http.StatusBadRequest, errOfCommandHandler.Error())
	}

	if userID == "" {
		fmt.Printf("userController.SignIn ERROR -> There was an error while binding json - ERROR: %v\n", "Internal Server Error")
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Create a map to represent the JSON data
	requestData := map[string]string{"userId": userID}

	// Convert the map to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	body, errOfHelper := helpers.HttpPostHelper(configuration.BACKEND_URL+"/api/v1/alpha/jwt", bytes.NewBuffer(jsonData))

	if errOfHelper != nil {
		fmt.Printf("userController.SignIn ERROR -> There was an error while sending request to jwt - ERROR: %v\n", errOfHelper.Error())
		return fiber.NewError(http.StatusInternalServerError, errOfHelper.Error())
	}

	if body == nil {
		fmt.Printf("userController.SignIn ERROR -> There was an error while sending request to jwt - ERROR: %v\n", "body is nil")
		return fiber.NewError(http.StatusBadRequest, "Body is nil")
	}

	var data map[string]string

	errOfData := json.Unmarshal([]byte(string(body)), &data)
	if errOfData != nil {
		fmt.Printf("userController.SignIn ERROR -> There was an error while binding data - ERROR: %v\n:", errOfData.Error())
		return fiber.NewError(http.StatusInternalServerError, errOfData.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		map[string]interface{}{
			"message": "User Successfully Sign In",
			"response": map[string]interface{}{
				"accessToken":  data["accessToken"],
				"refreshToken": data["refreshToken"],
			},
		},
	)
}

// GetUser godoc
//
//	@Summary		This method used for get all users
//	@Description	get all users
//	@Tags			User
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} []response.UserResponse
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/user [get]
func (u *UserController) GetUser(ctx *fiber.Ctx) error {
	users, err := u.userQueryService.GetUser(ctx.UserContext())

	if err != nil {
		fmt.Printf("userController.GetUser ERROR -> There was an error while getting users - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.ToUserResponseList(users))
}

// GetUserById godoc
//
//	@Summary		This method get user by given id
//	@Description	get user by id
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"userId"
//
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} response.UserResponse
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/user/{userId} [get]
func (u *UserController) GetUserById(ctx *fiber.Ctx) error {
	user, err := u.userQueryService.GetUserById(ctx.UserContext(), ctx.Params("userId"))
	userCtx := ctx.UserContext().Value("user").(*middlewares.UserContext)

	fmt.Println(userCtx.UserID)

	if err != nil {
		fmt.Printf("userController.GetUserById ERROR -> There was an error while getting user - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.ToUserResponse(user))
}
