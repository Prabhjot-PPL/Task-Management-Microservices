package persistance

import (
	"context"
	"fmt"
	"user_service/src/internal/adaptors/ports"
	"user_service/src/internal/core/dto"
	"user_service/src/pkg/hashing"
	"user_service/src/pkg/logger"

	"github.com/google/uuid"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(d *Database) ports.UserRepository {
	return &UserRepo{db: d}
}

func (u *UserRepo) CreateUser(ctx context.Context, user dto.UserDetails) error {

	hpassword, e := hashing.HashPassword(user.Password)
	if e != nil {
		fmt.Print("\n")
		logger.Log.Error("Unable to hash password : ", e)
	}

	var userId uuid.UUID
	err := u.db.db.QueryRowContext(ctx, `
		INSERT INTO user_details (username, email, password)
		VALUES ($1, $2, $3)
		
		RETURNING id
	`, user.Username, user.Email, hpassword).Scan(&userId)

	if err != nil {
		fmt.Print("\n")
		logger.Log.Error("Error inserting user_details: ", err)
		return err
	}

	return nil
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (dto.UserDetails, error) {
	var user dto.UserDetails

	query := `SELECT id, username, email, password, created_at, updated_at FROM user_details WHERE email=$1`
	row := u.db.db.QueryRowContext(ctx, query, email)

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Print("\n")
		logger.Log.Error("Error getting user_details: ", err)
		return user, err
	}

	return user, nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, userID string) (dto.UserDetails, error) {
	var user dto.UserDetails

	query := `SELECT id, username, email, created_at, updated_at FROM user_details WHERE id=$1`
	row := u.db.db.QueryRowContext(ctx, query, userID)

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Print("\n")
		logger.Log.Error("Error getting user_details: ", err)
		return user, err
	}

	return user, nil
}

func (u *UserRepo) UpdateUserInfo(ctx context.Context, userId string, newData dto.UserDetails) (dto.UserDetails, error) {

	// Condition to not update user_id, password, created_at, updated_at
	if newData.Id != uuid.Nil {
		fmt.Print("\n")
		logger.Log.Error("Cannot update user_id")
		return dto.UserDetails{}, nil
	}

	if newData.Password != "" {
		fmt.Print("\n")
		logger.Log.Error("Cannot update password")
		return dto.UserDetails{}, nil
	}

	if !newData.CreatedAt.IsZero() {
		fmt.Print("\n")
		logger.Log.Error("Cannot update created_at")
		return dto.UserDetails{}, nil
	}

	if !newData.UpdatedAt.IsZero() {
		fmt.Print("\n")
		logger.Log.Error("Cannot update updated_at")
		return dto.UserDetails{}, nil
	}

	// Updating user info
	query1 := `UPDATE user_details SET username=$1, email=$2 WHERE id=$3`
	_, err := u.db.db.ExecContext(ctx, query1, newData.Username, newData.Email, userId)

	// Getting updated user info
	var updatedData dto.UserDetails

	query2 := `SELECT id, username, email, created_at, updated_at FROM user_details WHERE id=$1`
	row2 := u.db.db.QueryRowContext(ctx, query2, userId)

	err = row2.Scan(&updatedData.Id, &updatedData.Username, &updatedData.Email, &updatedData.CreatedAt, &updatedData.UpdatedAt)
	if err != nil {
		fmt.Print("\n")
		logger.Log.Error("Error fetching updated user_details: ", err)
		return dto.UserDetails{}, err
	}

	return updatedData, nil

}
