package utils

import (
	"github.com/SchunckLeonardo/go-expert-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func OpenDBConnection(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)
	err = db.AutoMigrate(&entity.Product{})
	assert.Nil(t, err)

	return db
}
