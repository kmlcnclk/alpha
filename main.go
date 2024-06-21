package main

import (
	"alpha.com/configuration"
	_ "alpha.com/docs"
	"alpha.com/internal/alpha.com/application/controller"
	"alpha.com/internal/alpha.com/application/controller/response"
	"alpha.com/internal/alpha.com/application/handler/businessAccount"
	"alpha.com/internal/alpha.com/application/handler/job"
	"alpha.com/internal/alpha.com/application/handler/jobApply"
	"alpha.com/internal/alpha.com/application/handler/jwt"
	"alpha.com/internal/alpha.com/application/handler/user"
	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/application/web"
	"alpha.com/internal/alpha.com/pkg/mongodb"
	"alpha.com/internal/alpha.com/pkg/server"
	"alpha.com/internal/alpha.com/pkg/server/services"
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

	// User Dependency injection
	userRepository := repository.NewUserRepository(mongoClient)
	userService := services.NewUserService()
	userQueryService := query.NewUserQueryService(userRepository)
	userCommandHandler := user.NewCommandHandler(userRepository, userService)
	userController := controller.NewUserController(userQueryService, userCommandHandler, customValidator)

	// Jwt Dependency injection
	jwtRepository := repository.NewJwtRepository(mongoClient)
	jwtService := services.NewJwtService()
	jwtQueryService := query.NewJwtQueryService(jwtRepository, userQueryService, jwtService)
	jwtCommandHandler := jwt.NewCommandHandler(jwtRepository, jwtService, userQueryService)
	jwtController := controller.NewJwtController(jwtQueryService, jwtCommandHandler, customValidator)

	// Business Account Dependency injection
	businessAccountRepository := repository.NewBusinessAccountRepository(mongoClient)
	businessAccountQueryService := query.NewBusinessAccountQueryService(businessAccountRepository)
	businessAccountCommandHandler := businessAccount.NewCommandHandler(businessAccountRepository)
	businessAccountController := controller.NewBusinessAccountController(businessAccountQueryService, businessAccountCommandHandler, customValidator)

	// Job Dependency injection
	jobRepository := repository.NewJobRepository(mongoClient)
	jobQueryService := query.NewJobQueryService(jobRepository)
	jobCommandHandler := job.NewCommandHandler(jobRepository, businessAccountQueryService)
	jobController := controller.NewJobController(jobQueryService, jobCommandHandler, customValidator)

	// Job Apply Dependency injection
	jobApplyRepository := repository.NewJobApplyRepository(mongoClient)
	jobApplyQueryService := query.NewJobApplyQueryService(jobApplyRepository)
	jobApplyCommandHandler := jobApply.NewCommandHandler(jobApplyRepository, jobQueryService, userQueryService)
	jobApplyController := controller.NewJobApplyController(jobApplyQueryService, jobApplyCommandHandler, customValidator)

	// Router initializing
	web.InitRouter(app, userController, jwtController, businessAccountController, jobController, jobApplyController)

	// Start server
	server.NewServer(app).StartHttpServer(mongoClient)
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
