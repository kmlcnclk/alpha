package request

import "alpha.com/internal/alpha.com/application/handler/user"

type UserSignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

func (req *UserSignInRequest) ToCommand() user.CommandSignIn {
	return user.CommandSignIn{
		Email:    req.Email,
		Password: req.Password,
	}
}
