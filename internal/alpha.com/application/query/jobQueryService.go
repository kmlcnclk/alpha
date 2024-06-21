package query

import (
	"context"
	"errors"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
)

type IJobQueryService interface {
	GetAllJobs(ctx context.Context) ([]*domain.Job, error)
	GetByIDAndBusinessAccountID(ctx context.Context, id, businessAccountID string) (*domain.Job, error)
}

type jobQueryService struct {
	jobRepository repository.IJobRepository
}

func NewJobQueryService(jobRepository repository.IJobRepository) IJobQueryService {
	return &jobQueryService{
		jobRepository: jobRepository,
	}
}

func (u *jobQueryService) GetAllJobs(ctx context.Context) ([]*domain.Job, error) {
	jobs, err := u.jobRepository.Get(ctx)

	if err != nil {
		return nil, err
	}

	if jobs == nil {
		return nil, errors.New("not found jobs")
	}

	return jobs, nil
}

func (u *jobQueryService) GetByIDAndBusinessAccountID(ctx context.Context, id, businessAccountID string) (*domain.Job, error) {
	job, err := u.jobRepository.GetByIDAndBusinessAccountID(ctx, id, businessAccountID)

	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("not found Job with given id: " + id)
	}

	return job, nil
}
