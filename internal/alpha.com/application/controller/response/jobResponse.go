package response

import (
	"time"

	"alpha.com/internal/alpha.com/domain"
)

type JobResponse struct {
	Id                string    `json:"_id"`
	BusinessAccountID string    `json:"businessAccountId"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Price             float32   `json:"price"`
	Category          string    `json:"category"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func ToJobResponse(job *domain.Job) JobResponse {
	return JobResponse{
		Id:                job.Id.Hex(),
		BusinessAccountID: job.BusinessAccountID.Hex(),
		Name:              job.Name,
		Description:       job.Description,
		Price:             job.Price,
		Category:          job.Category,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
	}
}

func ToJobResponseList(jobs []*domain.Job) []JobResponse {
	var response = make([]JobResponse, 0)

	for _, job := range jobs {
		response = append(response, ToJobResponse(job))
	}

	return response
}
