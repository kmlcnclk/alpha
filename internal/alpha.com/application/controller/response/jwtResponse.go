package response

import (
	"alpha.com/internal/alpha.com/domain"
)

type JwtResponse struct {
	Id           string `json:"_id"`
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func ToJwtResponse(jwt *domain.Jwt) JwtResponse {
	return JwtResponse{
		Id:           jwt.Id.Hex(),
		UserID:       jwt.UserID.Hex(),
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	}
}

func ToJwtResponseList(jwts []*domain.Jwt) []JwtResponse {
	var response = make([]JwtResponse, 0)

	for _, jwt := range jwts {
		response = append(response, ToJwtResponse(jwt))
	}

	return response
}
