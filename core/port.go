package core

import (
	"github.com/SawitProRecruitment/UserService/core/domain"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen --build_flags=--mod=mod -destination=../generated/mock/port.go -package=mock . UserService
type UserService interface {
	Register(data *domain.User) (*domain.User, error)
	Get(id string) (*domain.User, error)
	Patch(id string, data *domain.User) (*domain.User, error)
}

//go:generate mockgen --build_flags=--mod=mod -destination=../generated/mock/port.go -package=mock . AuthService
type AuthService interface {
	Login(req *domain.AuthCred) (*domain.AuthData, error)
	VerifyAuthHeader(authHeader string) (id string, err error)
}

//go:generate mockgen --build_flags=--mod=mod -destination=../generated/mock/port.go -package=mock . UserRepo
type UserRepo interface {
	CreateUser(data *domain.User) (*domain.User, error)
	Login(phone string, password string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	PatchUserByID(id string, data *domain.User) (*domain.User, error)
}
