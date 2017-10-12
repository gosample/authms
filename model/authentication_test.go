package model_test

import (
	"testing"

	"github.com/tomogoma/authms/model"
	testingH "github.com/tomogoma/authms/testing"
)

func TestNewAuthentication(t *testing.T) {
	tt := []struct {
		name   string
		db     *testingH.DBMock
		jwter  *testingH.JWTMock
		opts   []model.Option
		expErr bool
	}{
		{
			name:   "min valid deps",
			db:     &testingH.DBMock{},
			jwter:  &testingH.JWTMock{},
			expErr: false,
		},
		{
			name:   "nil db",
			db:     nil,
			jwter:  &testingH.JWTMock{},
			expErr: true,
		},
		{
			name:   "nil jwter",
			db:     &testingH.DBMock{},
			jwter:  nil,
			expErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			a, err := model.NewAuthentication(tc.db, tc.jwter, tc.opts...)
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got an error: %v", err)
			}
			if a == nil {
				t.Fatalf("Got nil *Authentication")
			}
		})
	}
}

func TestAuthentication_RegisterSelf(t *testing.T) {
	tt := []struct {
		name              string
		db                model.AuthStore
		jwter             *testingH.JWTMock
		opts              []model.Option
		loginType         string
		userType          string
		identifier        string
		secret            []byte
		expErr            bool
		expNotImplemented bool
	}{
		{
			name:       "successful username",
			db:         &testingH.DBMock{},
			jwter:      &testingH.JWTMock{},
			loginType:  model.LoginTypeUsername,
			userType:   model.UserTypeIndividual,
			identifier: "johndoe",
			secret:     []byte("12345678"),
			expErr:     false,
		},
		{
			name:       "successful phone",
			db:         &testingH.DBMock{},
			jwter:      &testingH.JWTMock{},
			loginType:  model.LoginTypePhone,
			userType:   model.UserTypeIndividual,
			identifier: "+254712345678",
			secret:     []byte("12345678"),
			expErr:     false,
		},
		{
			name:       "successful email",
			db:         &testingH.DBMock{},
			jwter:      &testingH.JWTMock{},
			loginType:  model.LoginTypeEmail,
			userType:   model.UserTypeIndividual,
			identifier: "johndoe@mailinator.co.ke",
			secret:     []byte("12345678"),
			expErr:     false,
		},
		{
			name:  "successful facebook",
			db:    &testingH.DBMock{},
			jwter: &testingH.JWTMock{},
			opts: []model.Option{
				model.WithFacebookCl(&testingH.FacebookMock{}),
			},
			loginType:  model.LoginTypeFacebook,
			userType:   model.UserTypeIndividual,
			identifier: "johndoe",
			secret:     []byte("12345678"),
			expErr:     false,
		},
		{
			name:              "facebook unavailable",
			db:                &testingH.DBMock{},
			jwter:             &testingH.JWTMock{},
			opts:              []model.Option{},
			loginType:         model.LoginTypeFacebook,
			userType:          model.UserTypeIndividual,
			identifier:        "johndoe",
			secret:            []byte("12345678"),
			expErr:            true,
			expNotImplemented: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			a := newAuthentication(t, tc.db, tc.jwter, tc.opts...)
			usr, err := a.RegisterSelf(tc.loginType, tc.userType, tc.identifier, tc.secret)
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				if tc.expNotImplemented != a.IsNotImplementedError(err) {
					t.Fatalf("Expected IsNotImplemented %t, got %v",
						tc.expNotImplemented, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if usr == nil {
				t.Fatalf("Received nil user")
			}
		})
	}
}

func newAuthentication(t *testing.T, d model.AuthStore, j model.JWTEr, opts ...model.Option) *model.Authentication {
	a, err := model.NewAuthentication(d, j, opts...)
	if err != nil {
		t.Fatalf("Error setting up: new authentication: %v", err)
	}
	return a
}
