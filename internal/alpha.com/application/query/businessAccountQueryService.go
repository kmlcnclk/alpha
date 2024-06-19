package query

import (
	"context"
	"errors"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
)

type IBusinessAccountQueryService interface {
	GetAllBusinessAccounts(ctx context.Context) ([]*domain.BusinessAccount, error)
}

type businessAccountQueryService struct {
	businessAccountRepository repository.IBusinessAccountRepository
}

func NewBusinessAccountQueryService(businessAccountRepository repository.IBusinessAccountRepository) IBusinessAccountQueryService {
	return &businessAccountQueryService{
		businessAccountRepository: businessAccountRepository,
	}
}

func (u *businessAccountQueryService) GetAllBusinessAccounts(ctx context.Context) ([]*domain.BusinessAccount, error) {
	businessAccounts, err := u.businessAccountRepository.Get(ctx)

	if err != nil {
		return nil, err
	}

	if businessAccounts == nil {
		return nil, errors.New("not found businessAccounts")
	}

	return businessAccounts, nil
}
