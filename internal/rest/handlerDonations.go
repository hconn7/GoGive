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

func (apiCfg *ApiConfig) HandlerCreateDonation(w http.ResponseWriter, r *http.Request) {

	type Params struct {
		Amount        string `json:"amount"`
		Email         string `json:"email"`
		NonProfitName string `json:"non_profit_name"`
	}
	params := Params{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't decode params", err)
		return
	}
	user, err := apiCfg.DbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "User doesn't exist", err)
		return
	}
	nonProf, err := apiCfg.DbQueries.GetNonProfByName(r.Context(), params.NonProfitName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, http.StatusBadRequest, "No Non-Profit Exists", err)
		}
		RespondWithError(w, http.StatusBadRequest, "Error finding non profit", err)
	}
	donAmount, err := auth.CheckDonation(params.Amount)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't verify amount", err)
		return
	}
	donation, err := apiCfg.DbQueries.CreateDonation(r.Context(), database.CreateDonationParams{
		Amount:        int32(donAmount),
		UserID:        user.ID,
		NonProfitName: nonProf.Name})
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Donation failed to create", err)
	}
	type Response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Amount      int32     `json:"amount"`
		UserID      uuid.UUID `json:"user_id"`
		NonProfitID uuid.UUID `json:"non_profit_id"`
	}
	RespondWithJson(w, 200, Response{
		ID:          donation.ID,
		CreatedAt:   donation.CreatedAt,
		UpdatedAt:   donation.UpdatedAt,
		Amount:      int32(donAmount),
		UserID:      donation.UserID,
		NonProfitID: nonProf.ID,
	})

}

func (apiCfg ApiConfig) HandlerGetDonations(w http.ResponseWriter, r *http.Request) {

}
