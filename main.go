package main

import (
	"alpha.com/configuration"
	_ "alpha.com/docs"
	"alpha.com/internal/alpha.com/application/controller"
	"alpha.com/internal/alpha.com/application/controller/response"
	"alpha.com/internal/alpha.com/application/handler/jwt"
	"alpha.com/internal/alpha.com/application/handler/user"
	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/application/web"
	"alpha.com/internal/alpha.com/pkg/mongodb"
	"alpha.com/internal/alpha.com/pkg/server"
	"alpha.com/internal/alpha.com/pkg/server/helpers/jwtHelper"
	"alpha.com/internal/alpha.com/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title			Alpha Fiber Rest Api
// @version		1.0
// @description	This is a sample swagger for alpha rest api
// @contact.name	Alpha
// @contact.email	alpha@gmail.com
func main() {
	// fiber framework http server
	app := fiber.New(
		fiber.Config{
			// Override default error handler
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {

				// Default status code
				statusCode := fiber.StatusInternalServerError

				// Retrieve the custom error from fiber's context if it exists
				var customError response.CustomError

				if e, ok := err.(*fiber.Error); ok {
					// Fiber error, use its status code and message
					statusCode = e.Code
					customError = response.CustomError{
						StatusCode: statusCode,
						Message:    e.Message,
					}
				} else {

					// Non-fiber error, use default status code and message
					customError = response.CustomError{
						StatusCode: statusCode,
						Message:    err.Error(),
					}
				}

				// Send custom error response
				return ctx.Status(customError.StatusCode).JSON(customError)
			},
		},
	)

	app.Use(recover.New())

	configureSwaggerUi(app)

	mongoClient := mongodb.ConnectMongoDB()

	// custom validator initializing
	customValidator := validation.NewCustomValidator(validator.New())

	userRepository := repository.NewUserRepository(mongoClient)
	userQueryService := query.NewUserQueryService(userRepository)
	userCommandHandler := user.NewCommandHandler(userRepository)
	userController := controller.NewUserController(userQueryService, userCommandHandler, customValidator)

	jwtRepository := repository.NewJwtRepository(mongoClient)
	jwtHelper := jwtHelper.NewJwtHelper()
	jwtQueryService := query.NewJwtQueryService(jwtRepository, userQueryService, jwtHelper)
	jwtCommandHandler := jwt.NewCommandHandler(jwtRepository, jwtHelper)
	jwtController := controller.NewJwtController(jwtQueryService, jwtCommandHandler, customValidator)

	// Router initializing
	web.InitRouter(app, userController, jwtController)

	// Start server
	server.NewServer(app).StartHttpServer()
}

func configureSwaggerUi(app *fiber.App) {
	if configuration.Env != "prod" {
		// Swagger injection
		app.Get("/swagger/*", swagger.HandlerDefault)

		// Root path to SwaggerUI redirection
		app.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusMovedPermanently).Redirect("/swagger/index.html")
		})
	}
}
