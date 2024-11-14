package entity

import (
	"github.com/SchunckLeonardo/go-expert-api/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	name := "John Doe"
	email := "j@j.com"
	password := "123456"

	user, err := NewUser(name, email, password)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)

	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)

	assert.IsType(t, entity.ID{}, user.ID)
}

func TestUser_ValidatePassword(t *testing.T) {
	password := "123456"

	user, err := NewUser("", "", password)
	assert.Nil(t, err)

	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("wrong_password"))

	assert.NotEqual(t, password, user.Password)
}
