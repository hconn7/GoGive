package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hconn7/GoGive/internal/auth"
	"github.com/hconn7/GoGive/internal/database"
)

func (apiCfg *ApiConfig) HandlerCreateNonProfit(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Email   string `json:"email"`
	}
	type Response struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Address   string    `json:"address"`
		Email     string    `json:"email"`
		OwnerID   uuid.UUID `json:"owner_id"`
		Owner     bool      `json:"is_owner"`
	}
	params := Params{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error decoding", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error getting header", err)
		return
	}
	fmt.Println("Extracted Token:", token)
	_, err = auth.ValidateJWT(token, apiCfg.TokenSecret)
	if err != nil {
		RespondWithError(w, 409, "Unauthorized!", err)
		return
	}

	user, err := apiCfg.DbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, 409, "User doesn't exist, please register an account before registering a non-profit!", err)
			return
		}
		RespondWithError(w, 500, "Error getting user from DB", err)
		return
	}
	nonProf, err := apiCfg.DbQueries.CreateNonProfit(r.Context(), database.CreateNonProfitParams{
		Name:    params.Name,
		Address: params.Address,
		Email:   params.Email,
		OwnerID: user.ID,
	})
	if err != nil {
		RespondWithError(w, 500, "Error creating non-profit", err)
		return
	}

	if err := apiCfg.DbQueries.UpdateIsOwner(r.Context(), user.ID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to update user params", err)
	}

	RespondWithJson(w, 200, Response{
		ID:        nonProf.ID,
		Name:      nonProf.Name,
		CreatedAt: nonProf.CreatedAt,
		UpdatedAt: nonProf.UpdatedAt,
		Address:   nonProf.Address,
		Email:     nonProf.Email,
		OwnerID:   user.ID,
		Owner:     user.IsOwner,
	})

}

func (apiCfg *ApiConfig) HandlerNonProfDEL(w http.ResponseWriter, r *http.Request) {
	if err := apiCfg.DbQueries.DeleteNonProfits(r.Context()); err != nil {
		RespondWithError(w, 500, "Issue deleting non_profits", err)
	}
	RespondWithJson(w, 200, "Deleted Non-Profts")
}
