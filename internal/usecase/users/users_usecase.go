package users

import (
	"context"
	"errors"
	"time"

	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/users"
	"github.com/yeremiaaryo/gotu-assignment/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -package=users -source=users_usecase.go -destination=users_usecase_mock_test.go
type usersRepository interface {
	GetUser(ctx context.Context, email string) (*users.Model, error)
	InsertUser(ctx context.Context, model users.Model) (*users.Model, error)
}

type usecase struct {
	usersRepository usersRepository
	cfg             *configs.Config
}

func New(usersRepository usersRepository, cfg *configs.Config) *usecase {
	return &usecase{usersRepository: usersRepository, cfg: cfg}
}

func (u *usecase) CreateUser(ctx context.Context, req users.CreateUserRequest) (*users.Model, error) {
	user, err := u.usersRepository.GetUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("email already exists")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	model := users.Model{
		Email:     req.Email,
		Password:  string(pass),
		CreatedAt: now,
		UpdatedAt: now,
	}

	return u.usersRepository.InsertUser(ctx, model)
}

func (u *usecase) Login(ctx context.Context, req users.LoginRequest) (string, error) {
	user, err := u.usersRepository.GetUser(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return jwt.CreateToken(user.ID, user.Email, u.cfg.Service.SecretKey)
}
