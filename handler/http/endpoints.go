package sawithttp

import (
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/cons"
	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) UserRegister(ctx echo.Context) error {
	req := generated.UserRegisterRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	u := &domain.User{
		FullName:    req.FullName,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	}

	data, err := h.userSvc.Register(u)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	resp := generated.UserRegisterResponse{
		Id:          data.ID,
		FullName:    data.FullName,
		PhoneNumber: data.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) UserLogin(ctx echo.Context) error {
	req := generated.UserLoginRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	cred := &domain.AuthCred{
		PhoneNumber:   req.PhoneNumber,
		EncryptedPass: req.Password,
	}

	data, err := h.authSvc.Login(cred)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	resp := generated.UserLoginResponse{
		Id:          data.ID,
		AccessToken: data.AccessToken,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) UserDetail(ctx echo.Context, id openapi_types.UUID) error {
	header := ctx.Request().Header
	authHeader := header.Get("Authorization")
	claimID, err := h.authSvc.VerifyAuthHeader(authHeader)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusForbidden,
			err.Error(),
		)
	}

	if claimID != id.String() {
		return echo.NewHTTPError(
			http.StatusForbidden,
			cons.ErrInvalidAuthorized.Error(),
		)
	}

	data, err := h.userSvc.Get(id.String())
	if err != nil {
		return echo.NewHTTPError(
			http.StatusForbidden,
			err.Error(),
		)
	}

	resp := generated.UserDetailResponse{
		Id:          data.ID,
		FullName:    data.FullName,
		PhoneNumber: data.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) UserPatch(ctx echo.Context, id uuid.UUID) error {
	header := ctx.Request().Header
	authHeader := header.Get("Authorization")
	claimID, err := h.authSvc.VerifyAuthHeader(authHeader)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusForbidden,
			err.Error(),
		)
	}

	if claimID != id.String() {
		return echo.NewHTTPError(
			http.StatusForbidden,
			err.Error(),
		)
	}

	req := generated.UserPatchRequest{}
	err = ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	u := &domain.User{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}
	data, err := h.userSvc.Patch(id.String(), u)
	if err != nil {
		if errors.Is(err, cons.ErrDataConflict) {
			return echo.NewHTTPError(
				http.StatusConflict,
				err.Error(),
			)
		}

		if _, is := err.(interface{ Unwrap() []error }); is {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				err,
			)
		}

		return echo.NewHTTPError(
			http.StatusForbidden,
			err.Error(),
		)
	}

	resp := generated.UserPatchResponse{
		Id:          data.ID,
		FullName:    data.FullName,
		PhoneNumber: data.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}
