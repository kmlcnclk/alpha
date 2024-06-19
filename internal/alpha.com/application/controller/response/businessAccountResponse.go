package response

import (
	"time"

	"alpha.com/internal/alpha.com/domain"
)

type BusinessAccountResponse struct {
	Id          string    `json:"_id"`
	UserID      string    `json:"userID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func ToBusinessAccountResponse(businessAccount *domain.BusinessAccount) BusinessAccountResponse {
	return BusinessAccountResponse{
		Id:          businessAccount.Id.Hex(),
		UserID:      businessAccount.UserID.Hex(),
		Name:        businessAccount.Name,
		Description: businessAccount.Description,
		CreatedAt:   businessAccount.CreatedAt,
		UpdatedAt:   businessAccount.UpdatedAt,
	}
}

func ToBusinessAccountResponseList(businessAccounts []*domain.BusinessAccount) []BusinessAccountResponse {
	var response = make([]BusinessAccountResponse, 0)

	for _, businessAccount := range businessAccounts {
		response = append(response, ToBusinessAccountResponse(businessAccount))
	}

	return response
}
