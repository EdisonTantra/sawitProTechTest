package authsvc_test

import (
	"log"
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/SawitProRecruitment/UserService/core/port"
	"github.com/SawitProRecruitment/UserService/core/service/authsvc"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

const (
	PrivateKeyPath = "../../../generated/cert/sawitapp"
	PublicKeyPath  = "../../../generated/cert/sawitapp.pub"
	ExpDuration    = 2 * time.Hour
)

type testcaseLogin struct {
	name          string
	phone         string
	pass          string
	mockFunc      func(repo *port.MockUserRepo)
	assertionFunc func(tokenData *domain.AuthData, err error)
}

func TestService_Login(t *testing.T) {
	testcases := []testcaseLogin{
		{
			name:  "success login",
			phone: "+6285156305136",
			pass:  "Passw0rd!",
			mockFunc: func(repo *port.MockUserRepo) {
				repo.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:          "1234-1234-1234-1234",
						FullName:    "Edison",
						PhoneNumber: "+6285156305136",
						LoginCount:  2,
					}, nil).
					Times(1)
			},
			assertionFunc: func(tokenData *domain.AuthData, err error) {
				Expect(*tokenData).To(
					MatchFields(IgnoreExtras, Fields{
						"ID":          Equal("1234-1234-1234-1234"),
						"AccessToken": Not(BeNil()),
					}))
				Expect(err).To(BeNil())
			},
		},
	}

	var (
		mockRepo *port.MockUserRepo
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo = port.NewMockUserRepo(mockCtrl)
			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			opts := authsvc.ServiceOpts{
				PrvKeyPath:       PrivateKeyPath,
				PubKeyPath:       PublicKeyPath,
				TokenExpDuration: ExpDuration,
			}

			svc, err := authsvc.New(opts, mockRepo)
			if err != nil {
				log.Fatal(err)
			}

			req := &domain.User{
				PhoneNumber: tc.phone,
				Password:    tc.pass,
			}

			tokenData, err := svc.Login(req)
			tc.assertionFunc(tokenData, err)
		})
	}
}

type testcaseVerify struct {
	name          string
	authHeader    string
	isValid       bool
	mockFunc      func(repo *port.MockUserRepo)
	assertionFunc func(id string, err error)
}

func TestService_VerifyAuthHeader(t *testing.T) {
	testcases := []testcaseVerify{
		{
			name:       "failed invalid token format",
			authHeader: "token",
			assertionFunc: func(id string, err error) {
				Expect(id).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name:       "failed invalid token type",
			authHeader: "bear token",
			assertionFunc: func(id string, err error) {
				Expect(id).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name:       "failed invalid token",
			authHeader: "Bearer xxx",
			assertionFunc: func(id string, err error) {
				Expect(id).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			},
		},
	}

	var (
		mockRepo *port.MockUserRepo
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo = port.NewMockUserRepo(mockCtrl)
			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			opts := authsvc.ServiceOpts{
				PrvKeyPath:       PrivateKeyPath,
				PubKeyPath:       PublicKeyPath,
				TokenExpDuration: ExpDuration,
			}

			svc, err := authsvc.New(opts, mockRepo)
			if err != nil {
				log.Fatal(err)
			}

			id, err := svc.VerifyAuthHeader(tc.authHeader)
			tc.assertionFunc(id, err)
		})
	}
}
