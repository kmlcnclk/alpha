package request

import "alpha.com/internal/alpha.com/application/handler/jwt"

type JwtCreateRequest struct {
	UserID string `json:"userID"`
}

func (req *JwtCreateRequest) ToCommand() jwt.Command {
	return jwt.Command{
		UserID: req.UserID,
	}
}
