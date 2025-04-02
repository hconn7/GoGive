package rest

import (
	"net/http"

	"github.com/hconn7/GoGive/internal/auth"
)

func (apiCfg *ApiConfig) AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, _ := auth.GetBearerToken(r.Header)
		_, err := auth.ValidateJWT(token, apiCfg.TokenSecret)
		if err != nil {
			RespondWithJson(w, 409, "Could not validate JWT")
		}
		next.ServeHTTP(w, r)

	})

}
