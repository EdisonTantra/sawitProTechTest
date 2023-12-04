package sawithttp_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/cons"

	"github.com/google/uuid"

	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/SawitProRecruitment/UserService/core/port"
	sawithttp "github.com/SawitProRecruitment/UserService/handler/http"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

type testcaseRegister struct {
	name          string
	reqBody       string
	mockFunc      func(userSvc *port.MockUserService, authSvc *port.MockAuthService)
	assertionFunc func(recorder *httptest.ResponseRecorder, err error)
}

func TestHandler_UserRegister(t *testing.T) {
	testcases := []testcaseRegister{
		{
			name: "bad request invalid request",
			reqBody: `{
				"full_name": "ed",
				"password": "Passw@rd123",
				"phone_number": "+6285156305150"
			}`,
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				userSvc.EXPECT().
					Register(gomock.Any()).
					Return(nil, errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusBadRequest))
			},
		},
		{
			name: "success user register",
			reqBody: `{
				"full_name": "Edison Tantra",
				"password": "Passw@rd123",
				"phone_number": "+6285156305150"
			}`,
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				userSvc.EXPECT().
					Register(gomock.Any()).
					Return(&domain.User{
						FullName:    "Edison Tantra",
						Password:    "Passw@rd123",
						PhoneNumber: "+6285156305150",
					}, nil).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(err).To(BeNil())
			},
		},
	}

	const URLPath = "/api/v1/users/register"
	var (
		mockUserSvc *port.MockUserService
		mockAuthSvc *port.MockAuthService
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserSvc = port.NewMockUserService(mockCtrl)
			mockAuthSvc = port.NewMockAuthService(mockCtrl)

			logger := logrus.New()
			logger.SetOutput(io.Discard)
			handler := sawithttp.NewHandler(logger, nil, mockUserSvc, mockAuthSvc)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, URLPath, strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.mockFunc != nil {
				tc.mockFunc(mockUserSvc, mockAuthSvc)
			}

			err := handler.UserRegister(c)
			tc.assertionFunc(rec, err)
		})
	}
}

type testcaseLogin struct {
	name          string
	reqBody       string
	mockFunc      func(userSvc *port.MockUserService, authSvc *port.MockAuthService)
	assertionFunc func(recorder *httptest.ResponseRecorder, err error)
}

func TestHandler_UserLogin(t *testing.T) {
	testcases := []testcaseLogin{
		{
			name: "bad request login",
			reqBody: `{
				"phone_number": "+6285156305150",
				"password": "Password123@"
			}`,
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					Login(gomock.Any()).
					Return(nil, errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusBadRequest))
			},
		},
		{
			name: "success login",
			reqBody: `{
				"phone_number": "+6285156305150",
				"password": "Password123@"
			}`,
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					Login(gomock.Any()).
					Return(&domain.AuthData{
						ID:          "1234-1234-1234-1234",
						AccessToken: "eysomethingtoken",
					}, nil).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(err).To(BeNil())
			},
		},
	}

	const URLPath = "/api/v1/users/login"
	var (
		mockUserSvc *port.MockUserService
		mockAuthSvc *port.MockAuthService
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserSvc = port.NewMockUserService(mockCtrl)
			mockAuthSvc = port.NewMockAuthService(mockCtrl)

			logger := logrus.New()
			logger.SetOutput(io.Discard)
			handler := sawithttp.NewHandler(logger, nil, mockUserSvc, mockAuthSvc)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, URLPath, strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.mockFunc != nil {
				tc.mockFunc(mockUserSvc, mockAuthSvc)
			}

			err := handler.UserLogin(c)
			tc.assertionFunc(rec, err)
		})
	}
}

type testcaseDetail struct {
	name          string
	id            string
	mockFunc      func(userSvc *port.MockUserService, authSvc *port.MockAuthService)
	assertionFunc func(recorder *httptest.ResponseRecorder, err error)
}

func TestHandler_UserDetail(t *testing.T) {
	testcases := []testcaseDetail{
		{
			name: "forbidden invalid token",
			id:   "1234-1234-1234-1234",
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("", errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusForbidden))
			},
		},
		{
			name: "forbidden id does not match",
			id:   "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("abcd-abcd-abcd-abcd", nil).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusForbidden))
			},
		},
		{
			name: "forbidden failed get user data",
			id:   "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("9ae8810c-7b28-4c4c-8dbc-ed43be3da208", nil).
					Times(1)

				userSvc.EXPECT().
					Get(gomock.Any()).
					Return(nil, errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusForbidden))
			},
		},
		{
			name: "success get user data",
			id:   "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("9ae8810c-7b28-4c4c-8dbc-ed43be3da208", nil).
					Times(1)

				userSvc.EXPECT().
					Get(gomock.Any()).
					Return(&domain.User{
						ID:          "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
						FullName:    "Edison Tantra",
						PhoneNumber: "+6285156305136",
					}, nil).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(err).To(BeNil())
			},
		},
	}

	const URLPath = "/api/v1/users/1234-1234-1234-1234"
	var (
		mockUserSvc *port.MockUserService
		mockAuthSvc *port.MockAuthService
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserSvc = port.NewMockUserService(mockCtrl)
			mockAuthSvc = port.NewMockAuthService(mockCtrl)

			logger := logrus.New()
			logger.SetOutput(io.Discard)
			handler := sawithttp.NewHandler(logger, nil, mockUserSvc, mockAuthSvc)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, URLPath, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.mockFunc != nil {
				tc.mockFunc(mockUserSvc, mockAuthSvc)
			}

			validID, _ := uuid.Parse(tc.id)
			err := handler.UserDetail(c, validID)
			tc.assertionFunc(rec, err)
		})
	}
}

type testcasePatch struct {
	name          string
	id            string
	reqBody       string
	mockFunc      func(userSvc *port.MockUserService, authSvc *port.MockAuthService)
	assertionFunc func(recorder *httptest.ResponseRecorder, err error)
}

func TestHandler_UserPatch(t *testing.T) {
	testcases := []testcasePatch{
		{
			name:    "forbidden invalid token",
			id:      "1234-1234-1234-1234",
			reqBody: "{}",
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("", errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusForbidden))
			},
		},
		{
			name:    "forbidden id does not match",
			id:      "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
			reqBody: "{}",
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("abcd-abcd-abcd-abcd", nil).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusForbidden))
			},
		},
		{
			name: "forbidden failed patch user data",
			id:   "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
			reqBody: `{
				"full_name": "edison tantra",
				"phone_number": "+628515630513611"
			}`,
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("9ae8810c-7b28-4c4c-8dbc-ed43be3da208", nil).
					Times(1)

				userSvc.EXPECT().
					Patch(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusForbidden))
			},
		},
		{
			name: "forbidden conflict when patch user data",
			id:   "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
			reqBody: `{
				"full_name": "edison tantra",
				"phone_number": "+628515630513611"
			}`,
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("9ae8810c-7b28-4c4c-8dbc-ed43be3da208", nil).
					Times(1)

				userSvc.EXPECT().
					Patch(gomock.Any(), gomock.Any()).
					Return(nil, cons.ErrDataConflict).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(err).To(HaveOccurred())
				e := err.(*echo.HTTPError)
				Expect(e.Code).To(Equal(http.StatusConflict))
			},
		},
		{
			name: "success patch user data",
			id:   "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
			mockFunc: func(userSvc *port.MockUserService, authSvc *port.MockAuthService) {
				authSvc.EXPECT().
					VerifyAuthHeader(gomock.Any()).
					Return("9ae8810c-7b28-4c4c-8dbc-ed43be3da208", nil).
					Times(1)

				userSvc.EXPECT().
					Patch(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:          "9ae8810c-7b28-4c4c-8dbc-ed43be3da208",
						FullName:    "Edison Tantra",
						PhoneNumber: "+6285156305136",
					}, nil).
					Times(1)
			},
			assertionFunc: func(recorder *httptest.ResponseRecorder, err error) {
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(err).To(BeNil())
			},
		},
	}

	const URLPath = "/api/v1/users/1234-1234-1234-1234"
	var (
		mockUserSvc *port.MockUserService
		mockAuthSvc *port.MockAuthService
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserSvc = port.NewMockUserService(mockCtrl)
			mockAuthSvc = port.NewMockAuthService(mockCtrl)

			logger := logrus.New()
			logger.SetOutput(io.Discard)
			handler := sawithttp.NewHandler(logger, nil, mockUserSvc, mockAuthSvc)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, URLPath, strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.mockFunc != nil {
				tc.mockFunc(mockUserSvc, mockAuthSvc)
			}

			validID, _ := uuid.Parse(tc.id)
			err := handler.UserPatch(c, validID)
			tc.assertionFunc(rec, err)
		})
	}
}
