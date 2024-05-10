package services

import (
	"context"

	"github.com/ilhamgepe/shops/internal/models"
	"github.com/ilhamgepe/shops/internal/repositories"
)

type UserService struct {
	ur repositories.UserRepository
}

func NewUserService(ur repositories.UserRepository) *UserService {
	return &UserService{ur: ur}
}

func (us *UserService) GetById(ctx context.Context, id int64) (models.User, error) {
	return us.ur.GetUserById(ctx, id)
}
