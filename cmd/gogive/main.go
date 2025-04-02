package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hconn7/GoGive/internal/database"
	"github.com/hconn7/GoGive/internal/rest"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type httpServer struct {
	handler http.Handler
	address string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}
	tokenSecret := os.Getenv("TOKEN_SECRET")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't load DB", err)
	}

	dbQueries := database.New(db)
	apiCfg := rest.ApiConfig{
		DbQueries:   dbQueries,
		TokenSecret: tokenSecret,
	}

	//Server init
	mux := http.NewServeMux()
	httpServ := httpServer{handler: mux, address: ":8080"}

	//Handlers

	//Token
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevokeRefreshToken)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerValidateRefreshToken)
	//User Handler
	mux.HandleFunc("POST /users/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /users/create", apiCfg.HandlerCreateUser)
	mux.HandleFunc("DELETE /users/reset", apiCfg.HandlerReset)
	//Donation
	mux.HandleFunc("GET /donations/{user_id}", apiCfg.HandlerGetDonations)
	mux.HandleFunc("POST /api/donations", apiCfg.HandlerCreateDonation)

	//Non-Prof
	mux.HandleFunc("POST /non-profits/create", apiCfg.HandlerCreateNonProfit)
	mux.HandleFunc("DELETE /non-profits/reset", apiCfg.HandlerNonProfDEL)

	fmt.Printf("Serving at%s ", httpServ.address)
	http.ListenAndServe(httpServ.address, httpServ.handler)

}
