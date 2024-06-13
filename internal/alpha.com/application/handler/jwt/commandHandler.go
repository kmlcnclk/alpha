package jwt

import (
	"context"
	"fmt"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"alpha.com/internal/alpha.com/pkg/server/helpers/jwtHelper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICommandHandler interface {
	Create(ctx context.Context, command Command) (string, string, error)
	Get(ctx context.Context) ([]*domain.Jwt, error)
}

type commandHandler struct {
	jwtRepository repository.IJwtRepository
	jwtHelper     jwtHelper.IJwtHelper
}

func NewCommandHandler(jwtRepository repository.IJwtRepository, jwtHelper jwtHelper.IJwtHelper) ICommandHandler {
	return &commandHandler{
		jwtRepository: jwtRepository,
		jwtHelper:     jwtHelper,
	}
}

func (c *commandHandler) Create(ctx context.Context, command Command) (string, string, error) {
	accessToken, refreshToken, err := c.jwtHelper.CreateTokens(command.UserID)

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

	fmt.Printf("commandHandler.Create SUCCESS -> Jwt Tokens successfully created")

	return accessToken, refreshToken, nil
}

func (c *commandHandler) Get(ctx context.Context) ([]*domain.Jwt, error) {
	return c.jwtRepository.Get(ctx)
}

func (c *commandHandler) BuildEntity(accessToken, refreshToken string, userID primitive.ObjectID) *domain.Jwt {
	return &domain.Jwt{
		// Id:           uuid.NewString(),
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
