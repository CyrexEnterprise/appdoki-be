package app

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func (a *Application) AuthRouter(router *mux.Router) {
	authHandler := NewAuthHandler(a.conf.AppConfig, a.usersRepository)

	router.
		Methods(http.MethodGet).
		Path("/auth/login").
		HandlerFunc(authHandler.Login)

	router.
		Methods(http.MethodGet).
		Path("/auth/token").
		HandlerFunc(authHandler.Token)

	router.
		Methods(http.MethodGet).
		Path("/auth/url").
		HandlerFunc(authHandler.GetURL)
}

func (a *Application) JwtVerify(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const headerPrefix = "Bearer "
		tokenHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(tokenHeader, headerPrefix) {
			respondJSON(w, nil, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(tokenHeader, headerPrefix)

		verifier := a.conf.AppConfig.OIDCProvider.Verifier(&oidc.Config{
			ClientID: a.conf.AppConfig.GoogleOauth.ClientID,
		})

		parsedToken, err := verifier.Verify(r.Context(), token)
		if err != nil {
			respondJSON(w, nil, http.StatusUnauthorized)
			return
		}

		var idTokenClaims struct {
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		if err := parsedToken.Claims(&idTokenClaims); err != nil {
			respondInternalError(w)
			return
		}

		newReq := r.WithContext(context.WithValue(r.Context(), "user", idTokenClaims))

		next.ServeHTTP(w, newReq)
	}
}
