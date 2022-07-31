package api

import (
	"encoding/json"
	"net/http"
	"porcupine/pkg/config"
	"porcupine/pkg/models"
	"porcupine/pkg/services"
	"time"

	"github.com/upper/db/v4"
)

type Handlers struct {
	UserService  services.UserService
	TokenService services.TokenService
	CryptService services.CryptService
}

func (h *Handlers) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	userId := user.ID.Hex()
	pairId := "asdasdasdasvgdasdjhbaj"

	duration, _ := time.ParseDuration(config.ACCESS_TOKEN_AGE)

	accessToken := h.TokenService.SignJwt(pairId, userId, config.ACCESS_TOKEN_SECRET, time.Now().Add(duration).Unix())

	duration, _ = time.ParseDuration(config.REFRESH_TOKEN_AGE)
	expiration := time.Now().Add(duration).Unix()
	refreshToken := h.TokenService.SignJwt(pairId, userId, config.REFRESH_TOKEN_SECRET, expiration)

	err := h.TokenService.PushToStore(userId, pairId, time.Duration(expiration))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response := ErrorResponse{Message: err.Error(), Success: false}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	type Response struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	response := Response{Success: true, Data: map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}}
	json.NewEncoder(w).Encode(response)
	return
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Username string `json:"username" binding:"required"`
		Email    string ` json:"email" binding:"required"`
		Password string ` json:"password" binding:"required"`
	}

	var body ReqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response := ErrorResponse{Message: err.Error(), Success: false}
		json.NewEncoder(w).Encode(response)
		return
	}

	if validationErr := validate.Struct(&body); validationErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response := ErrorResponse{Message: validationErr.Error(), Success: false}
		json.NewEncoder(w).Encode(response)
		return
	}

	exist, _ := h.UserService.FindIfExists(db.Cond{
		"username": body.Username,
		"email":    body.Email,
	})
	if exist {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response := ErrorResponse{Message: "Username or email already Taken", Success: false}
		json.NewEncoder(w).Encode(response)
		return
	}

	hash, _ := h.CryptService.CryprtPassword(body.Password)
	user, err := h.UserService.Create(body.Username, body.Email, hash)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		response := ErrorResponse{Message: "Username or email already Taken", Success: false}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(user)
	return
}

// func (h *Handlers) BlackListTokens() {

// }

// func (h *Handlers) Me() {

// }

// func (h *Handlers) Rotate() {

// }
