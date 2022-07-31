package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/upper/db/v4"
)

var validate = validator.New()

func (h *Handlers) VerifyCredentials(next http.Handler) http.Handler {
	type ReqBody struct {
		Username string `json:"username" binding:"required"`
		Email    string ` json:"email" binding:"required"`
		Password string ` json:"password" binding:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		user, err := h.UserService.FindOne(db.Cond{"email": body.Email})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			response := ErrorResponse{Message: err.Error(), Success: false}
			json.NewEncoder(w).Encode(response)
			return
		}

		if !h.CryptService.ComparePasswordHash(user.Hash, body.Password) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			response := ErrorResponse{Message: "Invalid password credential", Success: false}
			json.NewEncoder(w).Encode(response)
			return
		}

		ctx := context.WithValue(r.Context(), "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
