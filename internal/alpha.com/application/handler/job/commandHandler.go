package job

import (
	"context"
	"fmt"
	"time"

	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command, userID string) error
}

type commandHandler struct {
	jobRepository               repository.IJobRepository
	businessAccountQueryService query.IBusinessAccountQueryService
}

func NewCommandHandler(jobRepository repository.IJobRepository, businessAccountQueryService query.IBusinessAccountQueryService) ICommandHandler {
	return &commandHandler{
		jobRepository:               jobRepository,
		businessAccountQueryService: businessAccountQueryService,
	}
}

func (c *commandHandler) Save(ctx context.Context, command Command, userID string) error {

	_, err := c.businessAccountQueryService.GetByIDAndUserID(ctx, command.BusinessAccountID, userID)

	if err != nil {
		fmt.Printf("commandHandler.Save ERROR -> Error was happened while finding Business Account with given id: %v Error:  %s\n", command.BusinessAccountID, err.Error())
		return err
	}

	businessAccountID, err := primitive.ObjectIDFromHex(command.BusinessAccountID)
	if err != nil {
		fmt.Printf("commandHandler.Save ERROR :  %s\n", err.Error())
		return err
	}

	newJob := c.BuildEntity(command, businessAccountID)

	err = c.jobRepository.Upsert(ctx, newJob)

	if err != nil {
		return err
	}

	return nil
}

func (c *commandHandler) BuildEntity(command Command, businessAccountID primitive.ObjectID) *domain.Job {
	return &domain.Job{
		BusinessAccountID: businessAccountID,
		Name:              command.Name,
		Description:       command.Description,
		Price:             command.Price,
		Category:          command.Category,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}
