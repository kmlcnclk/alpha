package controller

import (
	"fmt"
	"net/http"

	"alpha.com/internal/alpha.com/application/controller/request"
	"alpha.com/internal/alpha.com/application/controller/response"
	"alpha.com/internal/alpha.com/application/handler/job"
	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/pkg/server/middlewares"
	"alpha.com/internal/alpha.com/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type IJobController interface {
	Save(ctx *fiber.Ctx) error
	GetAllJobs(ctx *fiber.Ctx) error
}

type JobController struct {
	jobQueryService   query.IJobQueryService
	jobCommandHandler job.ICommandHandler
	customValidator   validation.ICustomValidator
}

func NewJobController(
	jobQueryService query.IJobQueryService,
	jobCommandHandler job.ICommandHandler,
	customValidator validation.ICustomValidator,
) IJobController {
	return &JobController{
		jobQueryService:   jobQueryService,
		jobCommandHandler: jobCommandHandler,
		customValidator:   customValidator,
	}
}

// Save godoc

//	@Summary		This method used for saving new job
//	@Description	saving new job
//
//	@Tags			Jobs
//	@Accept			json
//	@Produce		json
//
// @Param requestBody body request.JobCreateRequest nil "Handle Request Body"
//
// @Param Authorization header string true "Bearer {token}"
//
// @Success 200
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/job [post]
func (u *JobController) Save(ctx *fiber.Ctx) error {
	var req request.JobCreateRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		fmt.Printf("JobController.Save ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("JobController.Save STARTED with request: %#v\n", req)

	if err := u.customValidator.Validate(req); err != nil {
		fmt.Printf("JobController.Save INVALID request: %#v - ERROR: %#v\n", req, err)
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	userCtx := ctx.UserContext().Value("user").(*middlewares.UserContext)

	errOfCommandHandler := u.jobCommandHandler.Save(ctx.UserContext(), req.ToCommand(), userCtx.UserID)

	if errOfCommandHandler != nil {
		return fiber.NewError(http.StatusBadRequest, errOfCommandHandler.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		map[string]interface{}{
			"message": "Job Successfully Created",
		},
	)
}

// GetUser godoc
//
//	@Summary		This method used for get all jobs
//	@Description	get all jobs
//	@Tags			Jobs
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} []response.JobResponse
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/job [get]
func (u *JobController) GetAllJobs(ctx *fiber.Ctx) error {
	jobs, err := u.jobQueryService.GetAllJobs(ctx.UserContext())

	if err != nil {
		fmt.Printf("jobController.GetAllJobs ERROR -> There was an error while getting jobs - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.ToJobResponseList(jobs))
}
