package services

import "golang.org/x/crypto/bcrypt"

type IUserService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type userService struct {
}

func NewUserService() IUserService {
	return &userService{}
}

func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *userService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
