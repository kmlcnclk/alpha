package businessAccount

import (
	"context"
	"fmt"
	"time"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command, UserID string) error
}

type commandHandler struct {
	businessAccountRepository repository.IBusinessAccountRepository
}

func NewCommandHandler(businessAccountRepository repository.IBusinessAccountRepository) ICommandHandler {
	return &commandHandler{
		businessAccountRepository: businessAccountRepository,
	}
}

func (c *commandHandler) Save(ctx context.Context, command Command, UserID string) error {
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		fmt.Printf("commandHandler.Create ERROR :  %s\n", err.Error())
	}

	newBusinessAccount := c.BuildEntity(command, userID)

	err = c.businessAccountRepository.Upsert(ctx, newBusinessAccount)

	if err != nil {
		return err
	}

	return nil
}

func (c *commandHandler) BuildEntity(command Command, userID primitive.ObjectID) *domain.BusinessAccount {
	return &domain.BusinessAccount{
		UserID:      userID,
		Name:        command.Name,
		Description: command.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
