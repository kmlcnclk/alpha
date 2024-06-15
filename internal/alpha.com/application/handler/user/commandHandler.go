package user

import (
	"context"
	"fmt"
	"time"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"alpha.com/internal/alpha.com/pkg/server/services"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command) (string, error)
	SignIn(ctx context.Context, command CommandSignIn) (string, error)
}

type commandHandler struct {
	userRepository repository.IUserRepository
	userService    services.IUserService
}

func NewCommandHandler(userRepository repository.IUserRepository, userService services.IUserService) ICommandHandler {
	return &commandHandler{
		userRepository: userRepository,
		userService:    userService,
	}
}

func (c *commandHandler) SignIn(ctx context.Context, command CommandSignIn) (string, error) {
	user, err := c.userRepository.GetByEmail(ctx, command.Email)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("user could not find with given email: %s", command.Email)
	}

	if !c.userService.CheckPasswordHash(command.Password, user.Password) {
		return "", fmt.Errorf("password is not correct for given email: %s", command.Email)
	}

	return user.Id.Hex(), nil
}

func (c *commandHandler) Save(ctx context.Context, command Command) (string, error) {
	user, err := c.userRepository.GetByEmail(ctx, command.Email)

	if err != nil {
		return "", err
	}

	if user != nil {
		return "", fmt.Errorf("user Already Exist for given email: %s", command.Email)
	}

	hashedPassword, err := c.userService.HashPassword(command.Password)

	if err != nil {
		return "", fmt.Errorf("password could not hash: %s", err.Error())
	}

	newUser := c.BuildEntity(command, hashedPassword)

	objectID, err := c.userRepository.Upsert(ctx, newUser)

	if err != nil {
		return "", err
	}

	if objectID == "" {
		return "", fmt.Errorf("user could not be saved: %s", command.Email)
	}

	return objectID, nil
}

func (c *commandHandler) BuildEntity(command Command, hashedPassword string) *domain.User {
	return &domain.User{
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Email:     command.Email,
		Password:  hashedPassword,
		Age:       command.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
