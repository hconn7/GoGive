package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hconn7/GoGive/internal/auth"
	"github.com/hconn7/GoGive/internal/database"
)

type User struct {
	ID          uuid.UUID     `json:"id"`
	Email       string        `json:"email"`
	Password    string        `json:"Password"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Owner       bool          `json:"owner"`
	NonProfitID uuid.NullUUID `json:"non_profit_id"`
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params Params
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, 400, "Bad request, Couldn't decode", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, 500, "Internal Hashing error", err)
		return
	}

	user, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashedPass})
	if err != nil {
		RespondWithError(w, 500, "Failed to create user", err)
		return
	}

	RespondWithJson(w, http.StatusCreated, User{
		ID:          user.ID,
		Email:       params.Email,
		Password:    params.Password,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Owner:       user.IsOwner,
		NonProfitID: user.NonProfitID,
	})

}

func (apiCfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	params := Params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to decoode", err)
		return
	}

	user, err := apiCfg.DbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, http.StatusBadRequest, "User doesn't exist", err)
			return
		} else {

			RespondWithError(w, http.StatusBadRequest, "Couldnt not get user", err)
			return
		}
	}

	if err := auth.ComparePasswordHash(user.HashedPassword, params.Password); err != nil {
		RespondWithError(w, 409, "Wrong Password", err)
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()
	_, err = apiCfg.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{UserID: user.ID, Token: refreshToken})
	if err != nil {
		RespondWithError(w, 500, "Could not make RefreshToken in DB", err)
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, apiCfg.TokenSecret, time.Hour)
	if err != nil {
		RespondWithError(w, 500, "Could not make JWT", err)
		return
	}

	RespondWithJson(w, 200, Response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        jwtToken,
		RefreshToken: refreshToken,
	})
}
