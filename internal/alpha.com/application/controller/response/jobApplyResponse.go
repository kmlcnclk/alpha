package response

import (
	"time"

	"alpha.com/internal/alpha.com/domain"
)

type JobApplyResponse struct {
	Id        string    `json:"_id"`
	JobID     string    `json:"jobId"`
	UserID    string    `json:"userID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToJobApplyResponse(jobApply *domain.JobApply) JobApplyResponse {
	return JobApplyResponse{
		Id:        jobApply.Id.Hex(),
		JobID:     jobApply.JobID.Hex(),
		UserID:    jobApply.UserID.Hex(),
		CreatedAt: jobApply.CreatedAt,
		UpdatedAt: jobApply.UpdatedAt,
	}
}

func ToJobApplyResponseList(jobApplies []*domain.JobApply) []JobApplyResponse {
	var response = make([]JobApplyResponse, 0)

	for _, jobApply := range jobApplies {
		response = append(response, ToJobApplyResponse(jobApply))
	}

	return response
}
