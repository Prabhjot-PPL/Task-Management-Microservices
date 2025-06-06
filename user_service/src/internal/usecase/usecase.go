package usecase

import (
	"context"
	"user_service/src/internal/adaptors/ports"
	"user_service/src/internal/core/dto"
	"user_service/src/pkg"
	"user_service/src/pkg/logger"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) Service {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) RegisterUser(ctx context.Context, userData dto.UserDetails) error {
	err := u.userRepo.CreateUser(ctx, userData)
	return err
}

func (u *UserService) LoginUser(ctx context.Context, userRequestData dto.UserDetails) (dto.UserDetails, error) {

	dbUser, err := u.userRepo.GetUserByEmail(ctx, userRequestData.Email)
	if err != nil {
		print("\n")
		logger.Log.Error(err)
		return dto.UserDetails{}, err
	}

	err = pkg.CheckPassword(dbUser.Password, userRequestData.Password)
	if err != nil {
		print("\n")
		logger.Log.Error(err)
		return dto.UserDetails{}, err
	}

	return dbUser, nil

}

func (u *UserService) GetUserProfile(ctx context.Context, userID string) (dto.UserDetails, error) {
	return u.userRepo.GetUserByID(ctx, userID)
}

func (u *UserService) UpdateUserProfile(ctx context.Context, userId string, reqData dto.UserDetails) (dto.UserDetails, error) {
	return u.userRepo.UpdateUserInfo(ctx, userId, reqData)
}
