package controller

import (
	"fmt"
	"net/http"

	"alpha.com/internal/alpha.com/application/controller/request"
	"alpha.com/internal/alpha.com/application/controller/response"
	"alpha.com/internal/alpha.com/application/handler/jobApply"
	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type IJobApplyController interface {
	Save(ctx *fiber.Ctx) error
	GetAllJobApplies(ctx *fiber.Ctx) error
}

type JobApplyController struct {
	jobApplyQueryService   query.IJobApplyQueryService
	jobApplyCommandHandler jobApply.ICommandHandler
	customValidator        validation.ICustomValidator
}

func NewJobApplyController(
	jobApplyQueryService query.IJobApplyQueryService,
	jobApplyCommandHandler jobApply.ICommandHandler,
	customValidator validation.ICustomValidator,
) IJobApplyController {
	return &JobApplyController{
		jobApplyQueryService:   jobApplyQueryService,
		jobApplyCommandHandler: jobApplyCommandHandler,
		customValidator:        customValidator,
	}
}

// Save godoc

//	@Summary		This method used for saving new jobApply
//	@Description	saving new jobApply
//
//	@Tags			Job Applies
//	@Accept			json
//	@Produce		json
//
// @Param requestBody body request.JobApplyCreateRequest nil "Handle Request Body"
//
// @Param Authorization header string true "Bearer {token}"
//
// @Success 200
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/job-apply [post]
func (u *JobApplyController) Save(ctx *fiber.Ctx) error {
	var req request.JobApplyCreateRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		fmt.Printf("JobApplyController.Save ERROR -> There was an error while binding json - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("JobApplyController.Save STARTED with request: %#v\n", req)

	if err := u.customValidator.Validate(req); err != nil {
		fmt.Printf("JobApplyController.Save INVALID request: %#v - ERROR: %#v\n", req, err)
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	errOfCommandHandler := u.jobApplyCommandHandler.Save(ctx.UserContext(), req.ToCommand())

	if errOfCommandHandler != nil {
		return fiber.NewError(http.StatusBadRequest, errOfCommandHandler.Error())
	}

	return ctx.Status(http.StatusOK).JSON(
		map[string]interface{}{
			"message": "Job Apply Successfully Created",
		},
	)
}

// GetUser godoc
//
//	@Summary		This method used for get all job applies
//	@Description	get all job applies
//	@Tags			Job Applies
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} []response.JobApplyResponse
//
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/alpha/job-apply [get]
func (u *JobApplyController) GetAllJobApplies(ctx *fiber.Ctx) error {
	jobApplies, err := u.jobApplyQueryService.GetAllJobApplies(ctx.UserContext())

	if err != nil {
		fmt.Printf("jobApplyController.GetAllJobApplies ERROR -> There was an error while getting jobApplys - ERROR: %v\n", err.Error())
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.ToJobApplyResponseList(jobApplies))
}
