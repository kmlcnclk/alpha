package web

import (
	"fmt"
	"net/http"

	"alpha.com/internal/alpha.com/application/controller"
	"alpha.com/internal/alpha.com/pkg/server/middlewares"
	"github.com/gofiber/fiber/v2"
)

func InitRouter(app *fiber.App,
	userController controller.IUserController,
	jwtController controller.IJwtController,
	businessAccountController controller.IBusinessAccountController,
	jobController controller.IJobController,
	jobApplyController controller.IJobApplyController,
) {

	app.Get("/healthcheck", func(context *fiber.Ctx) error {
		fmt.Printf("Request sent to '/healthcheck' route -> Status: %v\n", http.StatusOK)
		return context.SendStatus(http.StatusOK)
	})

	alphaRouteGroup := app.Group("/api/v1/alpha")

	alphaRouteGroup.Get("/user", userController.GetUser)
	alphaRouteGroup.Post("/user", middlewares.IsEmailFormatCorrect, userController.Save)
	alphaRouteGroup.Post("/user/sign-in", middlewares.IsEmailFormatCorrect, userController.SignIn)
	alphaRouteGroup.Get("/user/:userId", middlewares.JwtMiddleware, userController.GetUserById)

	alphaRouteGroup.Post("/jwt", jwtController.Create)
	alphaRouteGroup.Post("/jwt/refresh", jwtController.Refresh)
	alphaRouteGroup.Get("/jwt", jwtController.GetJwt)

	alphaRouteGroup.Post("/business-account", middlewares.JwtMiddleware, businessAccountController.Save)
	alphaRouteGroup.Get("/business-account", businessAccountController.GetAllBusinessAccounts)

	alphaRouteGroup.Post("/job", middlewares.JwtMiddleware, jobController.Save)
	alphaRouteGroup.Get("/job", jobController.GetAllJobs)

	alphaRouteGroup.Post("/job-apply", middlewares.JwtMiddleware, jobApplyController.Save)
	alphaRouteGroup.Get("/job-apply", jobApplyController.GetAllJobApplies)
}
