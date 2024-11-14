package database

import (
	"github.com/SchunckLeonardo/go-expert-api/internal/entity"
	"github.com/SchunckLeonardo/go-expert-api/test/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestUser_Create(t *testing.T) {
	db := utils.OpenDBConnection(t)

	user, _ := entity.NewUser("John Doe", "j@j.com", "123456")
	userDB := NewUser(db)

	err := userDB.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotEqual(t, "123456", userFound.Password)

	err = db.Delete(&entity.User{}, "id = ?", user.ID).Error
	assert.Nil(t, err)

	err = db.First(&userFound, "id = ?", user.ID).Error
	log.Println(err)
	assert.NotNil(t, err)
}

func TestUser_FindByEmail(t *testing.T) {
	db := utils.OpenDBConnection(t)

	email := "j@j.com"

	user, _ := entity.NewUser("John Doe", email, "123456")
	userDB := NewUser(db)

	_ = userDB.Create(user)

	userFounded, err := userDB.FindByEmail(email)
	assert.Nil(t, err)
	assert.NotNil(t, userFounded)

	assert.Equal(t, user.ID, userFounded.ID)
	assert.Equal(t, user.Name, userFounded.Name)
	assert.Equal(t, user.Email, userFounded.Email)
	assert.NotNil(t, userFounded.Password)

	err = db.Delete(&entity.User{}, "id = ?", user.ID).Error
	assert.Nil(t, err)

	err = db.First(&entity.User{}, "id = ?", user.ID).Error
	log.Println(err)
	assert.NotNil(t, err)
}
