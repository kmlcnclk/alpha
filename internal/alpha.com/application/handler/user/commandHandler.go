package user

import (
	"context"
	"errors"
	"fmt"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"alpha.com/internal/alpha.com/pkg/server/helpers/userHelper"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command) (string, error)
	SignIn(ctx context.Context, command CommandSignIn) (string, error)
}

type commandHandler struct {
	userRepository repository.IUserRepository
}

func NewCommandHandler(userRepository repository.IUserRepository) ICommandHandler {
	return &commandHandler{
		userRepository: userRepository,
	}
}

func (c *commandHandler) SignIn(ctx context.Context, command CommandSignIn) (string, error) {
	user, err := c.userRepository.GetByEmail(ctx, command.Email)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New(fmt.Sprintf("User could not find with given email: %s", command.Email))
	}

	if !userHelper.CheckPasswordHash(command.Password, user.Password) {
		return "", errors.New(fmt.Sprintf("Password is not correct for given email: %s", command.Email))
	}

	return user.Id.Hex(), nil
}

func (c *commandHandler) Save(ctx context.Context, command Command) (string, error) {
	user, err := c.userRepository.GetByEmail(ctx, command.Email)

	if err != nil {
		return "", err
	}

	if user != nil {
		return "", errors.New(fmt.Sprintf("User Already Exist for given email: %s", command.Email))
	}

	hashedPassword, err := userHelper.HashPassword(command.Password)

	if err != nil {
		return "", errors.New(fmt.Sprintf("Password could not hash: %s", err.Error()))
	}

	newUser := c.BuildEntity(command, hashedPassword)

	objectID, err := c.userRepository.Upsert(ctx, newUser)

	if err != nil {
		return "", err
	}

	if objectID == "" {
		return "", errors.New(fmt.Sprintf("User could not be saved: %s", command.Email))
	}

	return objectID, nil
}

func (c *commandHandler) BuildEntity(command Command, hashedPassword string) *domain.User {
	return &domain.User{
		// Id:        uuid.NewString(),
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Email:     command.Email,
		Password:  hashedPassword,
		Age:       command.Age,
	}
}
