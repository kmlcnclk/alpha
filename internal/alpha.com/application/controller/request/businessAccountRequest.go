package request

import "alpha.com/internal/alpha.com/application/handler/businessAccount"

type BusinessAccountCreateRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description" validate:"required"`
}

func (req *BusinessAccountCreateRequest) ToCommand() businessAccount.Command {
	return businessAccount.Command{
		Name:        req.Name,
		Description: req.Description,
	}
}
