package adapters

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

//GetToken tries to extract Bearer token from the Authorization header
func GetToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")

	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}

	return ""
}

func extractUserFromToken(r *http.Request) string {
	token := GetToken(r)

	if token != "" {
		var p jwt.Parser
		claims := jwt.StandardClaims{}

		if _, _, err := p.ParseUnverified(token, &claims); err == nil && claims.Subject != "" {
			return claims.Subject
		}
	}
	return "-"
}

//Authenticate verifies the JWT token (if provided in the Authorization header)
func Authenticate(secret string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := GetToken(r)

			token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err == nil && token.Valid {
				h.ServeHTTP(w, r)
				return
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					fmt.Println("Invalid token")
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					fmt.Println("Token expired")
				} else {
					fmt.Println("Couldn't handle this token:", err)
				}
			} else {
				fmt.Println("Couldn't handle this token:", err)
			}

			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		})
	}
}
