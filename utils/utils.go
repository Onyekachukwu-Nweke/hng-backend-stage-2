package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/go-playground/validator/v10"
)

var jwtKey = []byte("your_jwt_secret_key")

type Claims struct {
	Email  string `json:"email"`
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateJWT(email, userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		Email:  email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Hour * 24)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func RespondWithValidationError(w http.ResponseWriter, err error) {
	var errors []ErrorResponse
	for _, err := range err.(validator.ValidationErrors) {
		var element ErrorResponse
		element.Field = err.StructNamespace()
		element.Message = err.Tag()
		errors = append(errors, element)
	}

	response := map[string]interface{}{
		"errors": errors,
	}
	RespondWithJSON(w, http.StatusUnprocessableEntity, response)
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
