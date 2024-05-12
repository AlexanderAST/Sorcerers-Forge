package model_test

import (
	model2 "Diploma/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model2.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model2.User {
				return model2.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "emptyEmail",
			u: func() *model2.User {
				u := model2.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := model2.TestUser(t)

	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
