package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"alpha.com/internal/alpha.com/application/query"
	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"alpha.com/internal/alpha.com/pkg/server/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICommandHandler interface {
	Create(ctx context.Context, command Command) (string, string, error)
	Refresh(ctx context.Context, userID string) (string, error)
	Update(ctx context.Context, userID, accessToken, refreshToken string) error
}

type commandHandler struct {
	jwtRepository    repository.IJwtRepository
	jwtService       services.IJwtService
	userQueryService query.IUserQueryService
}

func NewCommandHandler(jwtRepository repository.IJwtRepository, jwtService services.IJwtService, userQueryService query.IUserQueryService) ICommandHandler {
	return &commandHandler{
		jwtRepository:    jwtRepository,
		jwtService:       jwtService,
		userQueryService: userQueryService,
	}
}

func (c *commandHandler) Create(ctx context.Context, command Command) (string, string, error) {
	_, err := c.userQueryService.GetUserById(ctx, command.UserID)

	if err != nil {
		fmt.Printf("commandHandler.Create ERROR -> There was an error while finding user with given id: %v Error: %v\n", command.UserID, err.Error())
		return "", "", err
	}

	accessToken, refreshToken, err := c.jwtService.CreateTokens(command.UserID)

	if err != nil {

		fmt.Printf("commandHandler.Create ERROR -> There was an error while creating jwt tokens - ERROR: %v\n", err.Error())

		return "", "", err
	}

	userID, err := primitive.ObjectIDFromHex(command.UserID)
	if err != nil {
		fmt.Printf("commandHandler.Create ERROR :  %s\n", err.Error())
	}

	data := c.BuildEntity(accessToken, refreshToken, userID)

	if err = c.jwtRepository.Upsert(ctx, data); err != nil {

		fmt.Printf("commandHandler.Create ERROR -> There was an error while saving jwt tokens - ERROR: %v\n", err.Error())
		return "", "", err
	}

	fmt.Println("commandHandler.Create SUCCESS -> Jwt Tokens successfully created")

	return accessToken, refreshToken, nil
}

func (c *commandHandler) Refresh(ctx context.Context, userID string) (string, error) {
	_, err := c.userQueryService.GetUserById(ctx, userID)

	if err != nil {
		fmt.Printf("commandHandler.Refresh ERROR -> There was an error while finding user with given id: %v Error: %v\n", userID, err.Error())
		return "", err
	}

	accessToken, err := c.jwtService.CreateAccessToken(userID)

	if err != nil {
		fmt.Printf("commandHandler.Refresh ERROR -> There was an error while creating access token - ERROR: %v\n", err.Error())
		return "", err
	}

	if accessToken == "" {
		fmt.Println("commandHandler.Refresh ERROR -> Access Token is empty")
		return "", errors.New("access could not created")
	}

	return accessToken, nil
}

func (c *commandHandler) Update(ctx context.Context, userID, accessToken, refreshToken string) error {
	err := c.jwtRepository.Update(ctx, userID, accessToken, refreshToken)

	if err != nil {
		fmt.Printf("commandHandler.Update ERROR -> There was an error while updating jwt tokens - ERROR: %v\n", err.Error())
		return err
	}

	return nil
}

func (c *commandHandler) BuildEntity(accessToken, refreshToken string, userID primitive.ObjectID) *domain.Jwt {
	return &domain.Jwt{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
