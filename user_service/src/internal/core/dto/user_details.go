package dto

import (
	"time"

	// uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"github.com/google/uuid"
)

type UserDetails struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
