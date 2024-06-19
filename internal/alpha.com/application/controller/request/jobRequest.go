package request

import "alpha.com/internal/alpha.com/application/handler/job"

type JobCreateRequest struct {
	BusinessAccountID string  `json:"businessAccountId" validate:"required"`
	Name              string  `json:"name" validate:"required,min=2"`
	Description       string  `json:"description" validate:"required"`
	Price             float32 `json:"price" validate:"required"`
	Category          string  `json:"category" validate:"required"`
}

func (req *JobCreateRequest) ToCommand() job.Command {
	return job.Command{
		BusinessAccountID: req.BusinessAccountID,
		Name:              req.Name,
		Description:       req.Description,
		Price:             req.Price,
		Category:          req.Category,
	}
}
