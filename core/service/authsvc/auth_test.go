package authsvc_test

import (
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/SawitProRecruitment/UserService/core/service/authsvc"
	"github.com/SawitProRecruitment/UserService/generated/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

type testcaseLogin struct {
	name          string
	phone         string
	pass          string
	mockFunc      func(repo *mock.MockUserRepo)
	assertionFunc func(tokenData *domain.AuthData, err error)
}

func TestService_Login(t *testing.T) {
	testcases := []testcaseLogin{
		{
			name:  "failed login",
			phone: "+6285156305136",
			pass:  "Passw0rd!",
			mockFunc: func(repo *mock.MockUserRepo) {
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
				Expect(tokenData).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
	}

	var (
		mockRepo *mock.MockUserRepo
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo = mock.NewMockUserRepo(mockCtrl)
			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			//TODO benerin
			opts := authsvc.ServiceOpts{
				PrvKeyPath:       "privatepath",
				PubKeyPath:       "publicPath",
				TokenExpDuration: 2 * time.Hour,
				EncryptSecretKey: "abcdef",
			}

			svc, err := authsvc.New(opts, mockRepo)
			if err != nil {
				panic(err)
			}

			req := &domain.AuthCred{
				PhoneNumber:   tc.phone,
				EncryptedPass: tc.pass,
			}

			tokenData, err := svc.Login(req)
			tc.assertionFunc(tokenData, err)
		})
	}
}

type testcaseVerify struct {
	name          string
	authHeader    string
	mockFunc      func(repo *mock.MockUserRepo)
	assertionFunc func(id string, err error)
}

func TestService_VerifyAuthHeader(t *testing.T) {
	//TODO benerin
	testcases := []testcaseVerify{
		{
			name:       "failed ",
			authHeader: "Bearer xxxx",
			assertionFunc: func(id string, err error) {
				Expect(id).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
	}

	var (
		mockRepo *mock.MockUserRepo
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Default = NewGomegaWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo = mock.NewMockUserRepo(mockCtrl)
			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			//TODO benerin
			opts := authsvc.ServiceOpts{
				PrvKeyPath:       "privatepath",
				PubKeyPath:       "publicPath",
				TokenExpDuration: 2 * time.Hour,
				EncryptSecretKey: "abcdef",
			}

			svc, err := authsvc.New(opts, mockRepo)
			if err != nil {
				panic(err)
			}

			id, err := svc.VerifyAuthHeader(tc.authHeader)
			tc.assertionFunc(id, err)
		})
	}
}
