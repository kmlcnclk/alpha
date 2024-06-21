package jobApply

import (
	"context"
	"fmt"
	"time"

	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"alpha.com/internal/alpha.com/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command) error
}

type commandHandler struct {
	jobApplyRepository repository.IJobApplyRepository
	jobQueryService    query.IJobQueryService
	userQueryService   query.IUserQueryService
}

func NewCommandHandler(jobApplyRepository repository.IJobApplyRepository,
	jobQueryService query.IJobQueryService,
	userQueryService query.IUserQueryService,
) ICommandHandler {
	return &commandHandler{
		jobApplyRepository: jobApplyRepository,
		jobQueryService:    jobQueryService,
		userQueryService:   userQueryService,
	}
}

func (c *commandHandler) Save(ctx context.Context, command Command) error {
	userCtx := ctx.Value("user").(*utils.UserContext)

	_, err := c.userQueryService.GetUserById(ctx, userCtx.UserID)

	if err != nil {
		fmt.Printf("commandHandler.Save ERROR -> Error was happened while finding User with given id: %v Error:  %s\n", userCtx.UserID, err.Error())
		return err
	}

	_, err = c.jobQueryService.GetByIDAndBusinessAccountID(ctx, command.JobID, command.BusinessAccountID)

	if err != nil {
		fmt.Printf("commandHandler.Save ERROR -> Error was happened while finding Job with given id: %v Error:  %s\n", command.JobID, err.Error())
		return err
	}

	jobID, err := primitive.ObjectIDFromHex(command.JobID)
	if err != nil {
		fmt.Printf("commandHandler.Save ERROR :  %s\n", err.Error())
		return err
	}

	userID, err := primitive.ObjectIDFromHex(userCtx.UserID)
	if err != nil {
		fmt.Printf("commandHandler.Save ERROR :  %s\n", err.Error())
		return err
	}

	newJobApply := c.BuildEntity(jobID, userID)

	err = c.jobApplyRepository.Upsert(ctx, newJobApply)

	if err != nil {
		return err
	}

	return nil
}

func (c *commandHandler) BuildEntity(jobID, userID primitive.ObjectID) *domain.JobApply {
	return &domain.JobApply{
		JobID:     jobID,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
