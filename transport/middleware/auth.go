// transport/middleware/auth.go

package middleware

import (
	"bootcamp-auth-microservice/infras"
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	HeaderAuthorization = "Authorization"
)

type Authentication struct {
	DB     *infras.Conn
	Secret []byte
}

func ProvideAuthentication(db *infras.Conn, secret []byte) *Authentication {
	return &Authentication{
		DB:     db,
		Secret: secret,
	}
}

// Middleware to verify the JWT token
func (a *Authentication) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if token.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.Secret, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Check if the token is valid
		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			http.Error(w, "Token invalid", http.StatusUnauthorized)
			return
		}

		// Add the token to the request context for use in protected endpoints
		ctx := context.WithValue(r.Context(), "token", token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) VerifyTeacherJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.Secret, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Check if the token is valid and contains the role "teacher"
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid || claims["role"] != "teacher" {
			http.Error(w, "Unauthorized. Only teachers are allowed", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
