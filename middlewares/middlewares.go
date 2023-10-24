package middlewares

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"icos/server/responses"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	base64EncodedPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgTGF4mKVEa+eWX0S/+EWIfkkqbLba5WuQ1KKGRQz+P56Y0WNRbgjNl0CObndffmixbpgp4kg5jKq78HoFFP7bj0jQSNC3P26K9xPolFXbAlNJe41VMdI7xOkOF0D9GCplEylGlUlCgpaBnbloI4WcbH+RQ6n6Qp6MmNE+/xC3OMMhgEBacbiGtIR71N/HcDYDUORE335sSRpkrHhMxk3eWgZdIyfX88n9UkI3CtgNGIGgF8/w7ZYF2XBmVuv5+QE9d5fM9pZKWQnzBnsMJy4Xc+qZrZMI45KCHIW/DSFVGSsGboiVHSNVOu3mNhPSjvJtIH/7lItCG6m5zvBAvNf8QIDAQAB"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logHttpCall(r.Method + " " + r.URL.String())
		next(w, r)
	}
}

func logHttpCall(format string, args ...interface{}) {
	fmt.Printf(time.Now().Format(time.RFC3339)+"  \x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func JWTValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		splitToken := strings.Split(tokenString, "Bearer")
		if len(splitToken) < 2 {
			fmt.Println(tokenString)
			err := errors.New("not authorized")
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		reqToken := splitToken[1]
		reqToken = strings.TrimSpace(reqToken)

		// example token to validate
		// tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ6MXdOWU9YVFRpd2RhQjJyS1NqNElrUG1YTDJHbXVTOVFRNERJeXpWbWlrIn0.eyJleHAiOjE2ODIxMTQ3MjgsImlhdCI6MTY4MjExNDQyOCwianRpIjoiZDE5NGZlMmEtZDUyOC00YjhhLWI3OWItODEwNDMxNWExNzIxIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9BcnFncmlmbyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJmOWExMTYxYi00YzhiLTRmM2MtODViYy0xYmM1MDE0MWFmNTYiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ3ZWItYXBwIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1hcnFncmlmbyIsIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJ3ZWItYXBwIjp7InJvbGVzIjpbInVtYV9wcm90ZWN0aW9uIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InByb2ZpbGUgZW1haWwiLCJjbGllbnRJZCI6IndlYi1hcHAiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImNsaWVudEhvc3QiOiIxNzIuMTcuMC4xIiwicHJlZmVycmVkX3VzZXJuYW1lIjoic2VydmljZS1hY2NvdW50LXdlYi1hcHAiLCJjbGllbnRBZGRyZXNzIjoiMTcyLjE3LjAuMSJ9.vqFWx5mMESEMww8E5t1J8ZmoCw1R9qv1qlgaYaG7FQcd8B_sN223cDYMoqJF5y-Xv9zaJ094fUmyDtJHv-ZTkxw3R9AtjG0cCjqMxgBn1X2irlNYEmR5ZX73YXDUxY6XuABLyTGdh00bEcaUIyFR1Pver2UDjMf2okcV1FgEd0Z_94j4pjqtcY0nbsWIKnLoVoor7QV6ytWRpMG25DvrSVxciaOpogOHlUhaWtTfMz-mvfFg64i_S6rIuT84APnVe6weAuj92YS6bUzBif_gcgNeMdLrJChxWdPMK9G5mDAgLOqUv-X5fPOw1arigInV0nCJmKV7LG6Yc1UlDHdmiA"

		// The base64-encoded public key downloaded from Keycloak.
		// The screenshot in the question shows the correct place to get it.
		// It's much longer than "P184S3h7Ugjl56l31qeJ4FKvmBB4iikc".
		// base64EncodedPublicKey := `replaced with the public key downloaded from Keycloak`
		publicKey, err := parseKeycloakRSAPublicKey(base64EncodedPublicKey)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// return the public key that is used to validate the token.
			return publicKey, nil
		})
		if err != nil {
			// fmt.Println("Error parsing or validating token:", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		if !token.Valid {
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		fmt.Println("Claims:", claims)
		next(w, r)
	}
}

func parseKeycloakRSAPublicKey(base64Encoded string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return nil, err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	return nil, fmt.Errorf("unexpected key type %T", publicKey)
}
