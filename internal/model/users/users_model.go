package users

import (
	"github.com/yeremiaaryo/gotu-assignment/internal/response"
)

type (
	// Model is the user model that is retrieved from DB
	Model struct {
		ID        int64  `db:"id" json:"id"`
		Email     string `db:"email" json:"email"`
		Password  string `db:"password" json:"-"`
		CreatedAt int64  `db:"created_at" json:"created_at"`
		UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	}
)

// All request struct go below this
type (
	CreateUserRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

// All response struct go below this
type (
	UserResponse struct {
		response.BaseResponse
		User *Model `json:"user"`
	}

	LoginResponse struct {
		response.BaseResponse
		Token string `json:"token"`
	}
)
