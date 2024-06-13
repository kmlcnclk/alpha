package query

import (
	"context"
	"errors"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
)

type IUserQueryService interface {
	GetUser(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, userId string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userQueryService struct {
	userRepository repository.IUserRepository
}

func NewUserQueryService(userRepository repository.IUserRepository) IUserQueryService {
	return &userQueryService{
		userRepository: userRepository,
	}
}

func (u *userQueryService) GetUser(ctx context.Context) ([]*domain.User, error) {

	users, err := u.userRepository.Get(ctx)

	if err != nil {
		return nil, err
	}

	if users == nil {
		return nil, errors.New("not found users")
	}

	return users, nil
}

func (u *userQueryService) GetUserById(ctx context.Context, userId string) (*domain.User, error) {
	user, err := u.userRepository.GetById(ctx, userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("not found error")
	}

	return user, nil
}

func (u *userQueryService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	user, err := u.userRepository.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("not found error")
	}

	return user, nil
}
