package handlers

import (
	"encoding/json"
	"github.com/SchunckLeonardo/go-expert-api/internal/dto"
	"github.com/SchunckLeonardo/go-expert-api/internal/entity"
	"github.com/SchunckLeonardo/go-expert-api/internal/infra/database"
	entity2 "github.com/SchunckLeonardo/go-expert-api/pkg/entity"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(userDB database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDB:       userDB,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

// Create user godoc
//
//	@Summary		Create user
//	@Description	Create user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.CreateUserInput	true	"user request"
//	@Success		201
//	@Failure		400	{object}	entity.Error
//	@Failure		500	{object}	entity.Error
//	@Router			/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userDTO dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&userDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	user, err := entity.NewUser(userDTO.Name, userDTO.Email, userDTO.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetJWT godoc
//
//	@Summary		Get a user JWT
//	@Description	Get a user JWT
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.GetJWTInput	true	"user credentials"
//	@Success		200		{object}	dto.GetJWTOutput
//	@Failure		400		{object}	entity.Error
//	@Failure		500		{object}	entity.Error
//	@Router			/sessions [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var userJwtDto dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&userJwtDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	user, err := h.UserDB.FindByEmail(userJwtDto.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	if !user.ValidatePassword(userJwtDto.Password) {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := entity2.Error{Message: err.Error()}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	payload := map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	}

	_, token, _ := h.Jwt.Encode(payload)

	accessToken := dto.GetJWTOutput{AccessToken: token}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(accessToken)
}
