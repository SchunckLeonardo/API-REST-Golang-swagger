package handlers

import (
	"encoding/json"
	"github.com/SchunckLeonardo/go-expert-api/internal/dto"
	"github.com/SchunckLeonardo/go-expert-api/internal/entity"
	"github.com/SchunckLeonardo/go-expert-api/internal/infra/database"
	entity2 "github.com/SchunckLeonardo/go-expert-api/pkg/entity"
	"github.com/go-chi/chi/v5"
	"math"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

// CreateProduct godoc
//
//	@Summary		Create a product
//	@Description	Create products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.CreateProductInput	true	"product request"
//	@Success		201
//	@Failure		400	{object}	entity.Error
//	@Failure		500	{object}	entity.Error
//	@Router			/products [post]
//	@Security		ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDTO dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&productDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	product, err := entity.NewProduct(productDTO.Name, productDTO.Description, productDTO.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = h.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
//
//	@Summary		Get a product by ID
//	@Description	Get a product by ID
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"product ID"	Format(uuid)
//	@Success		200	{object}	entity.Product
//	@Failure		400	{object}	entity.Error
//	@Failure		404	{object}	entity.Error
//	@Failure		500	{object}	entity.Error
//	@Router			/products/{id} [get]
//	@Security		ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	Update a product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"product ID"	Format(uuid)
//	@Param			request	body	dto.UpdateProductInput	true	"Product fields that can be changed"
//	@Success		200
//	@Failure		400	{object}	entity.Error
//	@Failure		404	{object}	entity.Error
//	@Failure		500	{object}	entity.Error
//	@Router			/products/{id} [put]
//	@Security		ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var productDTO dto.UpdateProductInput
	err = json.NewDecoder(r.Body).Decode(&productDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if productDTO.Name != "" {
		product.Name = productDTO.Name
	}

	if productDTO.Description != "" {
		product.Description = productDTO.Description
	}

	if productDTO.Price >= 1.0 {
		product.Price = productDTO.Price
	}

	err = h.ProductDB.Update(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete a product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"product ID"	Format(uuid)
//	@Success		200
//	@Failure		404	{object}	entity.Error
//	@Failure		500	{object}	entity.Error
//	@Router			/products/{id} [delete]
//	@Security		ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FetchProducts godoc
//
//	@Summary		List products
//	@Description	Get all products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			page	query		string	false	"page number"
//	@Param			limit	query		string	false	"amount items"
//	@Param			sort	query		string	false	"sort asc or desc"
//	@Success		200		{object}	dto.FetchProductsOutput
//	@Failure		500		{object}	entity.Error
//	@Router			/products [get]
//	@Security		ApiKeyAuth
func (h *ProductHandler) FetchProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 1
	}

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	count, err := h.ProductDB.GetProductsCount()

	var totalPages float64
	totalPages = float64(count) / float64(limitInt)
	totalPagesCeil := int(math.Ceil(totalPages))

	response := dto.FetchProductsOutput{
		Products:    products,
		ItemsAmount: len(products),
		TotalPages:  totalPagesCeil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
