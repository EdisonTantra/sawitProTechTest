package sawithttp

import (
	"github.com/SawitProRecruitment/UserService/core"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/sirupsen/logrus"
)

var _ generated.ServerInterface = (*Handler)(nil)

type Handler struct {
	logger  *logrus.Logger
	userSvc core.UserService
	authSvc core.AuthService
}

func NewHandler(
	logger *logrus.Logger,
	userSvc core.UserService,
	authSvc core.AuthService,
) *Handler {
	return &Handler{
		logger:  logger,
		userSvc: userSvc,
		authSvc: authSvc,
	}
}
