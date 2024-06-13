package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
// @Param requestBody body request.UserCreteRequest nil "Handle Request Body"
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
	var req request.UserCreteRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		fmt.Printf("userController.Save ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}

	fmt.Printf("userController.Save STARTED with request: %#v\n", req)

	if err := u.customValidator.Validate(req); err != nil {
		fmt.Printf("userController.Save INVALID request: %#v\n - ERROR: %#v", req, err)
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	userID, errOfCommandHandler := u.userCommandHandler.Save(ctx.UserContext(), req.ToCommand())

	if errOfCommandHandler != nil {

		return ctx.Status(http.StatusBadRequest).JSON(
			response.CustomError{
				ErrorName:  "Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    errOfCommandHandler.Error(),
			})
	}

	if userID == "" {
		fmt.Printf("userController.Save ERROR -> There was an error while binding json - ERROR: %v\n", "Internal Server Error")

		return ctx.Status(http.StatusBadRequest).JSON("Internal Server Error")
	}

	// Create a map to represent the JSON data
	requestData := map[string]string{"userId": userID}

	// Convert the map to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	body, errOfHelper := helpers.HttpPostHelper("http://localhost:8080/api/v1/alpha/jwt", bytes.NewBuffer(jsonData))

	if errOfHelper != nil {
		fmt.Printf("userController.Save ERROR -> There was an error while sending request to jwt - ERROR: %v\n", errOfHelper.Error())
		return errOfHelper
	}

	if body == nil {
		fmt.Printf("userController.Save ERROR -> There was an error while sending request to jwt - ERROR: %v\n", "body is nil")
		return ctx.Status(http.StatusBadRequest).JSON("body is nil")
	}

	var data map[string]string

	errOfData := json.Unmarshal([]byte(string(body)), &data)
	if errOfData != nil {
		fmt.Printf("userController.Save ERROR -> There was an error while binding data - ERROR: %v\n:", errOfData.Error())
		return errOfData
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
		return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
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
		return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.ToUserResponse(user))
}
