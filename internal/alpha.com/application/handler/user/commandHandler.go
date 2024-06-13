package user

import (
	"context"
	"errors"
	"fmt"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command) (string, error)
}

type commandHandler struct {
	userRepository repository.IUserRepository
}

func NewCommandHandler(userRepository repository.IUserRepository) ICommandHandler {
	return &commandHandler{
		userRepository: userRepository,
	}
}

func (c *commandHandler) Save(ctx context.Context, command Command) (string, error) {
	user, err := c.userRepository.GetByEmail(ctx, command.Email)

	if err != nil {
		return "", err
	}

	if user != nil {
		return "", errors.New(fmt.Sprintf("User Already Exist for given email: %s", command.Email))
	}

	newUser := c.BuildEntity(command)

	objectID, err := c.userRepository.Upsert(ctx, newUser)

	if err != nil {
		return "", err
	}

	if objectID == "" {
		return "", errors.New(fmt.Sprintf("User could not be saved: %s", command.Email))
	}

	return objectID, nil
}

func (c *commandHandler) BuildEntity(command Command) *domain.User {
	return &domain.User{
		// Id:        uuid.NewString(),
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Email:     command.Email,
		Password:  command.Password,
		Age:       command.Age,
	}
}
