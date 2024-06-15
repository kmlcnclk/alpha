package query

import (
	"context"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
	"alpha.com/internal/alpha.com/pkg/server/services"
)

type IJwtQueryService interface {
	GetUserById(ctx context.Context, userId string) (*domain.User, error)
	ParseRefreshToken(ctx context.Context, refreshToken string) (string, error)
	Get(ctx context.Context) ([]*domain.Jwt, error)
}

type jwtQueryService struct {
	jwtRepository    repository.IJwtRepository
	userQueryService IUserQueryService
	jwtService       services.IJwtService
}

func NewJwtQueryService(jwtRepository repository.IJwtRepository, userQueryService IUserQueryService, jwtService services.IJwtService) IJwtQueryService {
	return &jwtQueryService{
		jwtRepository:    jwtRepository,
		userQueryService: userQueryService,
		jwtService:       jwtService,
	}
}

func (u *jwtQueryService) GetUserById(ctx context.Context, userId string) (*domain.User, error) {
	user, err := u.userQueryService.GetUserById(ctx, userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (u *jwtQueryService) ParseRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userID, err := u.jwtService.ParseRefreshToken(refreshToken)

	if err != nil {
		return "", err
	}

	return userID, nil
}

func (c *jwtQueryService) Get(ctx context.Context) ([]*domain.Jwt, error) {
	return c.jwtRepository.Get(ctx)
}
