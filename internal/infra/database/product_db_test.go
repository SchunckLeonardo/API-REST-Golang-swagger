package database

import (
	"fmt"
	"github.com/SchunckLeonardo/go-expert-api/internal/entity"
	"github.com/SchunckLeonardo/go-expert-api/test/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"testing"
)

func TestProduct_Create(t *testing.T) {
	db := utils.OpenDBConnection(t)

	product, _ := entity.NewProduct("Product 1", "Description 1", 80.0)
	productDB := NewProduct(db)

	err := productDB.Create(product)
	assert.Nil(t, err)
	assert.NotEmpty(t, product.ID)

	err = db.Delete(&entity.Product{}, "id = ?", product.ID).Error
	assert.Nil(t, err)

	err = db.First(&entity.Product{}, "id = ?", product.ID).Error
	log.Println(err)
	assert.NotNil(t, err)
}

func TestProduct_FindAll(t *testing.T) {
	db := utils.OpenDBConnection(t)

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), fmt.Sprintf("Description %d", i), rand.Float64())
		assert.NoError(t, err)
		db.Create(product)
	}

	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)
}

func TestProduct_FindByID(t *testing.T) {
	db := utils.OpenDBConnection(t)

	product, err := entity.NewProduct("Laptop", "Macbook M1", 1100.00)
	assert.Nil(t, err)

	db.Create(&product)

	productDB := NewProduct(db)

	productFounded, err := productDB.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.NotNil(t, productFounded)

	assert.Equal(t, productFounded.ID, product.ID)
	assert.Equal(t, productFounded.Name, product.Name)
	assert.Equal(t, productFounded.Description, product.Description)
	assert.Equal(t, productFounded.Price, product.Price)
}

func TestProduct_Update(t *testing.T) {
	db := utils.OpenDBConnection(t)

	product, err := entity.NewProduct("Laptop", "Macbook M1", 1100.00)
	assert.Nil(t, err)

	db.Create(&product)

	productDB := NewProduct(db)

	product.Name = "Laptop 2"

	err = productDB.Update(product)
	assert.Nil(t, err)

	var productUpdated entity.Product
	err = db.Where("id = ?", product.ID).First(&productUpdated).Error
	assert.Nil(t, err)

	assert.Equal(t, productUpdated.Name, product.Name)
}

func TestProduct_Delete(t *testing.T) {
	db := utils.OpenDBConnection(t)

	product, err := entity.NewProduct("Laptop", "Macbook M1", 1100.00)
	assert.Nil(t, err)

	db.Create(&product)

	productDB := NewProduct(db)

	err = productDB.Delete(product.ID.String())
	assert.Nil(t, err)

	var productUpdated entity.Product
	err = db.Where("id = ?", product.ID).First(&productUpdated).Error
	assert.Error(t, err, gorm.ErrRecordNotFound)
}

func TestProduct_GetProductsCount(t *testing.T) {
	db := utils.OpenDBConnection(t)

	var products []entity.Product

	for i := 1; i < 5; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Laptop %d", i), "Macbook M1", 1100.00)
		assert.Nil(t, err)

		products = append(products, *product)

		db.Create(&product)
	}

	productDb := NewProduct(db)

	count, err := productDb.GetProductsCount()
	assert.Nil(t, err)
	assert.NotNil(t, count)

	assert.Equal(t, len(products), count)
}
