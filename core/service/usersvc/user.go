package usersvc

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/cons"
	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/SawitProRecruitment/UserService/core/port"
)

var _ port.UserService = (*Service)(nil)

type Service struct {
	repo port.UserRepo
}

func New(repo port.UserRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (svc *Service) Register(data *domain.User) (*domain.User, error) {
	data.FullName = strings.TrimSpace(data.FullName)
	data.PhoneNumber = strings.TrimSpace(data.PhoneNumber)
	data.Password = strings.TrimSpace(data.Password)

	if data.FullName == "" || data.PhoneNumber == "" || data.Password == "" {
		return nil, errors.New("fullname, phone and password required")
	}

	err := validateUserData(data)
	if err != nil {
		return nil, err
	}

	newUser, err := svc.repo.CreateUser(data)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (svc *Service) Get(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID required")
	}

	data, err := svc.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (svc *Service) Patch(id string, data *domain.User) (*domain.User, error) {
	data.FullName = strings.TrimSpace(data.FullName)
	data.PhoneNumber = strings.TrimSpace(data.PhoneNumber)

	err := validateUserData(data)
	if err != nil {
		return nil, err
	}

	res, err := svc.repo.PatchUserByID(id, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func validateUserData(data *domain.User) error {
	var errWrap error
	if data.FullName != "" {
		err := validateFullName(data.FullName)
		if err != nil {
			errWrap = err
		}
	}

	if data.Password != "" {
		err := validatePassword(data.Password)
		if err != nil {
			if errWrap == nil {
				errWrap = err
			} else {
				errWrap = fmt.Errorf("%w; %w", errWrap, err)
			}
		}
	}

	if data.PhoneNumber != "" {
		err := validatePhoneNumber(data.PhoneNumber)
		if err != nil {
			if errWrap == nil {
				errWrap = err
			} else {
				errWrap = fmt.Errorf("%w; %w", errWrap, err)
			}
		}
	}

	if errWrap != nil {
		return errWrap
	}

	return nil
}

func validateFullName(name string) error {
	if len(name) < cons.MinLengthName || len(name) > cons.MaxLengthName {
		return cons.ErrInvalidNameLength
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < cons.MinLengthPass || len(password) > cons.MaxLengthPass {
		return cons.ErrInvalidPasswordLength
	}

	var number, capital, symbol bool
	for _, r := range password {
		switch {
		case unicode.IsNumber(r):
			number = true
		case unicode.IsUpper(r):
			capital = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			symbol = true
		}
	}

	if !number || !capital || !symbol {
		return cons.ErrInvalidPasswordFormat
	}

	return nil
}

func validatePhoneNumber(phone string) error {
	if !strings.HasPrefix(phone, cons.PrefixPhoneID) {
		return cons.ErrInvalidPhonePrefix
	}

	cleanPhone := strings.Replace(phone, cons.PrefixPhoneID, "0", 1)
	if len(cleanPhone) < cons.MinLengthPhone || len(cleanPhone) > cons.MaxLengthPhone {
		return cons.ErrInvalidPhoneLength
	}

	return nil
}
