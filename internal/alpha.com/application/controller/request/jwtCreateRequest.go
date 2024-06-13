package request

import "alpha.com/internal/alpha.com/application/handler/jwt"

type JwtCreteRequest struct {
	UserID string `json:"userID"`
}

func (req *JwtCreteRequest) ToCommand() jwt.Command {
	return jwt.Command{
		UserID: req.UserID,
	}
}
