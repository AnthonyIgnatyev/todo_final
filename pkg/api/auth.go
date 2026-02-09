package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cfg "todo_final/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey string = "superSecretKey"

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonError(w, "Wrong method, need POST", http.StatusBadRequest)
		return
	}

	password := cfg.CfgStruct.Server.Password

	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password != password {
		writeJsonError(w, "Wrong password", http.StatusUnauthorized)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pass_hash": fmt.Sprintf("%x", sha256.Sum256([]byte(password))),
		"exp":       time.Now().Add(8 * time.Hour).Unix(),
	})

	signedToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token": signedToken,
	}

	w.Header().Set("Content-Type", "application/json")
	writeJson(w, response)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := cfg.CfgStruct.Server.Password
		if len(password) > 0 {
			var cookieJwt string

			cookie, err := r.Cookie("token")
			if err != nil {
				writeJsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cookieJwt = cookie.Value
			var valid bool

			token, err := jwt.Parse(cookieJwt, func(t *jwt.Token) (interface{}, error) {

				return []byte(secretKey), nil
			})
			if err != nil {
				fmt.Println("JWT parce error")
				writeJsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			valid = token.Valid

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
