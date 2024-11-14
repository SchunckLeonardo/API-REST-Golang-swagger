package entity

import (
	"github.com/SchunckLeonardo/go-expert-api/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProduct(t *testing.T) {
	name := "Product 1"
	description := "Product 1 description"
	price := 500.50

	product, err := NewProduct(name, description, price)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, name, product.Name)
	assert.Equal(t, description, product.Description)
	assert.Equal(t, price, product.Price)
	assert.NotEmpty(t, price, product.Price)
	assert.IsType(t, entity.ID{}, product.ID)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", "", 10.0)
	assert.Nil(t, product)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	product, err := NewProduct("Product 1", "", 0)
	assert.Nil(t, product)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	product, err := NewProduct("Product 1", "", -50)
	assert.Nil(t, product)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProduct_Validate(t *testing.T) {
	product, err := NewProduct("Product 1", "", 10)
	assert.Nil(t, err)
	assert.NotNil(t, product)

	assert.Nil(t, product.Validate())
}
