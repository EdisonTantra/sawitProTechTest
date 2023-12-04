package sawithttp

import (
	"github.com/SawitProRecruitment/UserService/core/port"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/lib/locker"
	"github.com/sirupsen/logrus"
)

var _ generated.ServerInterface = (*Handler)(nil)

type Handler struct {
	logger    *logrus.Logger
	libLocker *locker.Locker
	userSvc   port.UserService
	authSvc   port.AuthService
}

func NewHandler(
	logger *logrus.Logger,
	libLocker *locker.Locker,
	userSvc port.UserService,
	authSvc port.AuthService,
) *Handler {
	return &Handler{
		logger:    logger,
		libLocker: libLocker,
		userSvc:   userSvc,
		authSvc:   authSvc,
	}
}
