package usersvc_test

import (
	"errors"
	"testing"

	"github.com/SawitProRecruitment/UserService/cons"
	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/SawitProRecruitment/UserService/core/service/usersvc"
	"github.com/SawitProRecruitment/UserService/generated/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

type testcaseValidate struct {
	name string
	user *domain.User
	want error
}

func TestValidateFullName(t *testing.T) {
	testcases := []testcaseValidate{
		{
			name: "failed less than min length",
			user: &domain.User{
				FullName: "ed",
			},
			want: cons.ErrInvalidNameLength,
		},
		{
			name: "failed greater than max length",
			user: &domain.User{
				FullName: "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghija",
			},
			want: cons.ErrInvalidNameLength,
		},
		{
			name: "success with min length",
			user: &domain.User{
				FullName: "abc",
			},
			want: nil,
		},
		{
			name: "success with max length",
			user: &domain.User{
				FullName: "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij",
			},
			want: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := usersvc.ValidateFullName(tc.user.FullName)
			if !errors.Is(got, tc.want) {
				t.Errorf("error: expect %s not %s", tc.want, got)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	testcases := []testcaseValidate{
		{
			name: "failed less than min length",
			user: &domain.User{
				Password: "abcde",
			},
			want: cons.ErrInvalidPasswordLength,
		},
		{
			name: "failed greater than max length",
			user: &domain.User{
				Password: "Password12345@Password12Password12Password12Password12Password12a",
			},
			want: cons.ErrInvalidPasswordLength,
		},
		{
			name: "failed only lowercase",
			user: &domain.User{
				Password: "abcdef",
			},
			want: cons.ErrInvalidPasswordFormat,
		},
		{
			name: "failed only one uppercase",
			user: &domain.User{
				Password: "abcdeF",
			},
			want: cons.ErrInvalidPasswordFormat,
		},
		{
			name: "failed only one uppercase, one number",
			user: &domain.User{
				Password: "abcd4F",
			},
			want: cons.ErrInvalidPasswordFormat,
		},
		{
			name: "failed only one number",
			user: &domain.User{
				Password: "abcd4f",
			},
			want: cons.ErrInvalidPasswordFormat,
		},
		{
			name: "failed only one number, one symbol",
			user: &domain.User{
				Password: "abcd4$",
			},
			want: cons.ErrInvalidPasswordFormat,
		},
		{
			name: "failed only one symbol",
			user: &domain.User{
				Password: "abcd4$",
			},
			want: cons.ErrInvalidPasswordFormat,
		},
		{
			name: "failed only one uppercase, one symbol",
			user: &domain.User{
				Password: "abcdE$",
			},
			want: cons.ErrInvalidPasswordFormat,
		},
		{
			name: "success with min length",
			user: &domain.User{
				Password: "Edi$0n",
			},
			want: nil,
		},
		{
			name: "success with max length",
			user: &domain.User{
				Password: "Password12345@Password12Password12Password12Password12Password12",
			},
			want: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := usersvc.ValidatePassword(tc.user.Password)
			if !errors.Is(got, tc.want) {
				t.Errorf("error: expect %s not %s", tc.want, got)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	testcases := []testcaseValidate{
		{
			name: "failed dont have prefix +62",
			user: &domain.User{
				PhoneNumber: "085156305136",
			},
			want: cons.ErrInvalidPhonePrefix,
		},
		{
			name: "failed less than min length",
			user: &domain.User{
				PhoneNumber: "+6285156305",
			},
			want: cons.ErrInvalidPhoneLength,
		},
		{
			name: "failed greater than max length",
			user: &domain.User{
				PhoneNumber: "+628515630513612",
			},
			want: cons.ErrInvalidPhoneLength,
		},
		{
			name: "failed with min length",
			user: &domain.User{
				PhoneNumber: "+62851563051",
			},
			want: nil,
		},
		{
			name: "success with max length",
			user: &domain.User{
				PhoneNumber: "+62851563051361",
			},
			want: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := usersvc.ValidatePhoneNumber(tc.user.PhoneNumber)
			if !errors.Is(got, tc.want) {
				t.Errorf("error: expect %s not %s", tc.want, got)
			}
		})
	}
}

func TestValidateUserData(t *testing.T) {
	testcases := []testcaseValidate{
		{
			name: "failed invalid name length",
			user: &domain.User{
				FullName:    "Ed",
				PhoneNumber: "+6285156305136",
				Password:    "Passw0rd@",
			},
			want: cons.ErrInvalidNameLength,
		},
		{
			name: "failed invalid phone length",
			user: &domain.User{
				FullName:    "Edison",
				PhoneNumber: "+6285156305",
				Password:    "Passw0rd@",
			},
			want: cons.ErrInvalidPhoneLength,
		},
		{
			name: "failed invalid password length",
			user: &domain.User{
				FullName:    "Edison",
				PhoneNumber: "+6285156305136",
				Password:    "Pass",
			},
			want: cons.ErrInvalidPasswordLength,
		},
		{
			name: "failed invalid name and phone length",
			user: &domain.User{
				FullName:    "Ed",
				PhoneNumber: "+6285156305",
				Password:    "Passw0rd@",
			},
			want: errors.Join(cons.ErrInvalidNameLength, cons.ErrInvalidPhoneLength),
		},
		{
			name: "failed invalid name, phone and password length",
			user: &domain.User{
				FullName:    "Ed",
				PhoneNumber: "+6285156305",
				Password:    "Pass",
			},
			want: errors.Join(
				cons.ErrInvalidNameLength,
				cons.ErrInvalidPasswordLength,
				cons.ErrInvalidPhoneLength,
			),
		},
		{
			name: "success validate user",
			user: &domain.User{
				FullName:    "Edison",
				PhoneNumber: "+6285156305136",
				Password:    "Passw0rd@",
			},
			want: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := usersvc.ValidateUserData(tc.user)
			if got != nil || tc.want != nil {
				if got.Error() != tc.want.Error() {
					t.Errorf("error: expect %s not %s", tc.want, got)
				}
			} else {
				if !errors.Is(got, tc.want) {
					t.Errorf("error: expect %s not %s", tc.want, got)
				}
			}
		})
	}
}

type testcaseRegister struct {
	name          string
	user          *domain.User
	mockFunc      func(repo *mock.MockUserRepo)
	assertionFunc func(newUser *domain.User, err error)
}

func TestService_Register(t *testing.T) {
	testcases := []testcaseRegister{
		{
			name: "failed fullname empty",
			user: &domain.User{
				FullName:    "",
				Password:    "abcde",
				PhoneNumber: "+6212345",
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "failed phone number empty",
			user: &domain.User{
				FullName:    "Edison",
				Password:    "abcde",
				PhoneNumber: "",
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "failed password empty",
			user: &domain.User{
				FullName:    "Edison",
				Password:    "",
				PhoneNumber: "+625156305136",
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "failed validate data request",
			user: &domain.User{
				FullName:    "Edison",
				Password:    "password",
				PhoneNumber: "+625156305136",
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "failed repo create user",
			user: &domain.User{
				FullName:    "Edison",
				Password:    "Password123@",
				PhoneNumber: "+625156305136",
			},
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().
					CreateUser(gomock.Any()).
					Return(nil, errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "success register new user",
			user: &domain.User{
				FullName:    "Edison",
				Password:    "Password123@",
				PhoneNumber: "+625156305136",
			},
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().
					CreateUser(gomock.Any()).
					Return(&domain.User{
						ID:          "1234-1234-1234-1234",
						FullName:    "Edison",
						PhoneNumber: "+625156305136",
					}, nil).
					Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(*newUser).To(
					MatchFields(IgnoreExtras, Fields{
						"ID":          Equal("1234-1234-1234-1234"),
						"FullName":    Equal("Edison"),
						"PhoneNumber": Equal("+625156305136"),
					}))
				Expect(err).To(BeNil())
			},
		},
	}

	var (
		mockRepo *mock.MockUserRepo
		svc      *usersvc.Service
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
			svc = usersvc.New(mockRepo)

			newUser, err := svc.Register(tc.user)
			tc.assertionFunc(newUser, err)
		})
	}
}

type testcaseGet struct {
	name          string
	id            string
	mockFunc      func(repo *mock.MockUserRepo)
	assertionFunc func(newUser *domain.User, err error)
}

func TestService_Get(t *testing.T) {
	testcases := []testcaseGet{
		{
			name: "failed user ID empty",
			id:   "",
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "success get by user ID",
			id:   "1234-1234-1234-1234",
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&domain.User{
					ID:          "1234-1234-1234-1234",
					FullName:    "Edison",
					PhoneNumber: "+625156305136",
				}, nil).Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(*newUser).To(
					MatchFields(IgnoreExtras, Fields{
						"ID":          Equal("1234-1234-1234-1234"),
						"FullName":    Equal("Edison"),
						"PhoneNumber": Equal("+625156305136"),
					}))
				Expect(err).To(BeNil())
			},
		},
	}

	var (
		mockRepo *mock.MockUserRepo
		svc      *usersvc.Service
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
			svc = usersvc.New(mockRepo)

			newUser, err := svc.Get(tc.id)
			tc.assertionFunc(newUser, err)
		})
	}
}

type testcasePatch struct {
	name          string
	id            string
	user          *domain.User
	mockFunc      func(repo *mock.MockUserRepo)
	assertionFunc func(newUser *domain.User, err error)
}

func TestService_Patch(t *testing.T) {
	testcases := []testcasePatch{
		{
			name: "failed validate data",
			id:   "1234-1234-1234-1234",
			user: &domain.User{
				FullName:    "Edison",
				Password:    "Password123",
				PhoneNumber: "+625156305136",
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "failed patch user",
			id:   "1234-1234-1234-1234",
			user: &domain.User{
				FullName:    "Edison",
				Password:    "Password123@",
				PhoneNumber: "+625156305136",
			},
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().
					PatchUserByID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error occurred")).
					Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(newUser).To(BeNil())
				Expect(err).To(HaveOccurred())
			},
		},
		{
			name: "success only fullname",
			user: &domain.User{
				FullName: "Edison Tantra",
			},
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().
					PatchUserByID(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:          "1234-1234-1234-1234",
						FullName:    "Edison Tantra",
						Password:    "Password123@",
						PhoneNumber: "+625156305136",
					}, nil).
					Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(*newUser).To(
					MatchFields(IgnoreExtras, Fields{
						"ID":          Equal("1234-1234-1234-1234"),
						"FullName":    Equal("Edison Tantra"),
						"Password":    Equal("Password123@"),
						"PhoneNumber": Equal("+625156305136"),
					}))
				Expect(err).To(BeNil())
			},
		},
		{
			name: "success only password",
			user: &domain.User{
				Password: "Password12345!",
			},
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().
					PatchUserByID(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:          "1234-1234-1234-1234",
						FullName:    "Edison",
						Password:    "Password12345!",
						PhoneNumber: "+625156305136",
					}, nil).
					Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(*newUser).To(
					MatchFields(IgnoreExtras, Fields{
						"ID":          Equal("1234-1234-1234-1234"),
						"FullName":    Equal("Edison"),
						"Password":    Equal("Password12345!"),
						"PhoneNumber": Equal("+625156305136"),
					}))
				Expect(err).To(BeNil())
			},
		},
		{
			name: "success only phone number",
			user: &domain.User{
				PhoneNumber: "+62515630513678",
			},
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().
					PatchUserByID(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:          "1234-1234-1234-1234",
						FullName:    "Edison",
						Password:    "Password123@",
						PhoneNumber: "+62515630513678",
					}, nil).
					Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(*newUser).To(
					MatchFields(IgnoreExtras, Fields{
						"ID":          Equal("1234-1234-1234-1234"),
						"FullName":    Equal("Edison"),
						"Password":    Equal("Password123@"),
						"PhoneNumber": Equal("+62515630513678"),
					}))
				Expect(err).To(BeNil())
			},
		},
		{
			name: "success patch fullname, password, and phone",
			user: &domain.User{
				FullName:    "Edison Tantra",
				Password:    "Password12345!",
				PhoneNumber: "+62515630513678",
			},
			mockFunc: func(repo *mock.MockUserRepo) {
				repo.EXPECT().
					PatchUserByID(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:          "1234-1234-1234-1234",
						FullName:    "Edison Tantra",
						Password:    "Password12345!",
						PhoneNumber: "+62515630513678",
					}, nil).
					Times(1)
			},
			assertionFunc: func(newUser *domain.User, err error) {
				Expect(*newUser).To(
					MatchFields(IgnoreExtras, Fields{
						"ID":          Equal("1234-1234-1234-1234"),
						"FullName":    Equal("Edison Tantra"),
						"Password":    Equal("Password12345!"),
						"PhoneNumber": Equal("+62515630513678"),
					}))
				Expect(err).To(BeNil())
			},
		},
	}

	var (
		mockRepo *mock.MockUserRepo
		svc      *usersvc.Service
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
			svc = usersvc.New(mockRepo)

			newUser, err := svc.Patch(tc.id, tc.user)
			tc.assertionFunc(newUser, err)
		})
	}
}
