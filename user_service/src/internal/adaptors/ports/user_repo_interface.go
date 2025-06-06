package ports

import (
	"context"
	"user_service/src/internal/core/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userInfo dto.UserDetails) error
	GetUserByEmail(ctx context.Context, email string) (dto.UserDetails, error)
	GetUserByID(ctx context.Context, userID string) (dto.UserDetails, error)
	UpdateUserInfo(ctx context.Context, userId string, newData dto.UserDetails) (dto.UserDetails, error)
}
