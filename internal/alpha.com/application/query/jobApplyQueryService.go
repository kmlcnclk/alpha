package query

import (
	"context"
	"errors"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
)

type IJobApplyQueryService interface {
	GetAllJobApplies(ctx context.Context) ([]*domain.JobApply, error)
}

type jobApplyQueryService struct {
	jobApplyRepository repository.IJobApplyRepository
}

func NewJobApplyQueryService(jobApplyRepository repository.IJobApplyRepository) IJobApplyQueryService {
	return &jobApplyQueryService{
		jobApplyRepository: jobApplyRepository,
	}
}

func (u *jobApplyQueryService) GetAllJobApplies(ctx context.Context) ([]*domain.JobApply, error) {
	jobApplies, err := u.jobApplyRepository.Get(ctx)

	if err != nil {
		return nil, err
	}

	if jobApplies == nil {
		return nil, errors.New("not found job applies")
	}

	return jobApplies, nil
}
