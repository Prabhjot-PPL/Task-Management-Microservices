package usecase

import (
	"context"
	"user_service/src/internal/core/dto"
)

type Service interface {
	RegisterUser(ctx context.Context, userData dto.UserDetails) error
	LoginUser(ctx context.Context, userRequestData dto.UserDetails) (dto.UserDetails, error)
	GetUserProfile(ctx context.Context, userID string) (dto.UserDetails, error)
	UpdateUserProfile(ctx context.Context, userId string, reqData dto.UserDetails) (dto.UserDetails, error)
}
