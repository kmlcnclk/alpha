package query

import (
	"context"
	"errors"

	"alpha.com/internal/alpha.com/application/repository"
	"alpha.com/internal/alpha.com/domain"
)

type IJobQueryService interface {
	GetAllJobs(ctx context.Context) ([]*domain.Job, error)
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
