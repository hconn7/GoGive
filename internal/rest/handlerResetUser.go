package rest

import (
	"net/http"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if err := cfg.DbQueries.DeleteUsers(r.Context()); err != nil {
		RespondWithError(w, 500, "Couldn't reset users", err)
	}

	RespondWithJson(w, 200, "")
}
