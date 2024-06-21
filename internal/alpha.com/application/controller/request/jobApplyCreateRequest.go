package request

import "alpha.com/internal/alpha.com/application/handler/jobApply"

type JobApplyCreateRequest struct {
	JobID             string `json:"jobId" validate:"required"`
	BusinessAccountID string `json:"businessAccountId" validate:"required"`
}

func (req *JobApplyCreateRequest) ToCommand() jobApply.Command {
	return jobApply.Command{
		JobID:             req.JobID,
		BusinessAccountID: req.BusinessAccountID,
	}
}
