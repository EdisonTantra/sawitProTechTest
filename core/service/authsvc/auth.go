package authsvc

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/cons"
	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/SawitProRecruitment/UserService/core/port"
	"github.com/golang-jwt/jwt/v5"
)

var _ port.AuthService = (*Service)(nil)

type Service struct {
	tokenPrvKey      []byte
	tokenPubKey      []byte
	tokenExpDuration time.Duration
	repo             port.UserRepo
}

type ServiceOpts struct {
	PrvKeyPath       string
	PubKeyPath       string
	TokenExpDuration time.Duration
}

func New(opts ServiceOpts, repo port.UserRepo) (*Service, error) {
	prvKey, err := os.ReadFile(opts.PrvKeyPath)
	if err != nil {
		return nil, err
	}

	pubKey, err := os.ReadFile(opts.PubKeyPath)
	if err != nil {
		return nil, err
	}

	return &Service{
		tokenPrvKey:      prvKey,
		tokenPubKey:      pubKey,
		tokenExpDuration: opts.TokenExpDuration,
		repo:             repo,
	}, nil
}

func (svc *Service) Login(req *domain.User) (*domain.AuthData, error) {
	data, err := svc.repo.Login(req.PhoneNumber, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := svc.generateAccessToken(data)
	if err != nil {
		return nil, err
	}

	res := &domain.AuthData{
		ID:          data.ID,
		AccessToken: token,
	}

	return res, nil
}

func (svc *Service) VerifyAuthHeader(authHeader string) (string, error) {
	split := strings.Split(strings.TrimSpace(authHeader), " ")
	if len(split) != 2 {
		return "", cons.ErrInvalidToken
	}

	if split[0] != cons.AuthTokenType {
		return "", cons.ErrInvalidToken
	}

	token := split[1]
	id, err := svc.verifyToken(token)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (svc *Service) generateAccessToken(data *domain.User) (string, error) {
	jwtIssuer := "sawitApp"
	jwtAud := []string{
		"sawitApp",
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(svc.tokenPrvKey)
	if err != nil {
		return "", err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   data.ID,
		Audience:  jwtAud,
		Issuer:    jwtIssuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(svc.tokenExpDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (svc *Service) verifyToken(tokenStr string) (string, error) {
	var id string
	key, err := jwt.ParseRSAPublicKeyFromPEM(svc.tokenPubKey)
	if err != nil {
		return id, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return id, cons.ErrJWTSignMethod
		}

		return key, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return id, cons.ErrJWTFormat
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return id, cons.ErrJWTSign
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return id, cons.ErrJWTExpired
		default:
			return id, err
		}
	}

	if token.Valid {
		id, err = token.Claims.GetSubject()
		if err != nil {
			return id, err
		}
	}

	return id, nil
}
