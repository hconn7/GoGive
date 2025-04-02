package rest

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/hconn7/GoGive/internal/auth"
)

func (apiCfg *ApiConfig) HandlerValidateRefreshToken(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	refreshTok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "No Token\n", err)
		return
	}

	fullToken, err := apiCfg.DbQueries.GetRefreshTokenByToken(r.Context(), refreshTok)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Token not recognized", err)
		return
	}

	if time.Now().After(fullToken.ExpiresAt) || fullToken.RevokedAt.Valid {
		RespondWithError(w, 401, "Token Expired or revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(fullToken.UserID, apiCfg.TokenSecret, time.Hour)
	if err != nil {
		RespondWithError(w, 500, "Error making token", err)
		return
	}

	RespondWithJson(w, 200, Response{Token: accessToken})

}

func (apiCfg *ApiConfig) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "No Token\n", err)
		return
	}
	_, err = apiCfg.DbQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, 409, "Token doesn't exist", err)
		}
		RespondWithError(w, http.StatusUnauthorized, "Token not recognized", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
