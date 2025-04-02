package rest

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hconn7/GoGive/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

type UserT struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type ResponseLogin struct {
	UserT
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
type apiConfig struct {
	DbQueries   *database.Queries
	TokenSecret string
}

func TestUsers(t *testing.T) {
	if err := godotenv.Load("../../cmd/gogive/.env"); err != nil {
		log.Fatal("Error loading .env")
	}
	tokenSecret := os.Getenv("TOKEN_SECRET")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't load DB", err)
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{
		DbQueries:   dbQueries,
		TokenSecret: tokenSecret,
	}

	registerPayload := map[string]any{
		"email":    "test@example.com",
		"password": "password123",
	}

	request, err := sendPostRequest("/users/create", registerPayload, "")
	assert.NoError(t, err)
	assert.Equal(t, 201, request.StatusCode)

	requestLog, err := sendPostRequest("/users/login", registerPayload, "")
	assert.NoError(t, err)
	responseLog := ResponseLogin{}
	parseJSONResponse(requestLog, &responseLog)
	user, err := apiCfg.DbQueries.GetUserByEmail(context.Background(), responseLog.Email)
	if err != nil {
		t.Error("Error getting user")
	}
	assert.Equal(t, responseLog.ID, user.ID)

	token := responseLog.Token
	type Params struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Email   string `json:"email"`
	}

	reqeustProf, err := sendPostRequest("/non-profits/create", Params{
		Address: "121 coral drive",
		Name:    "henryConner",
		Email:   user.Email},
		token)
	type ResponseNonProf struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Address   string    `json:"address"`
		Email     string    `json:"email"`
		OwnerID   uuid.UUID `json:"owner_id"`
		Owner     bool      `json:"is_owner"`
	}
	respNonProf := ResponseNonProf{}
	parseJSONResponse(reqeustProf, &respNonProf)
	nonProf, err := apiCfg.DbQueries.GetNonProfByName(context.Background(), respNonProf.Name)
	if err != nil {
		t.Error("Could not fetch non prof")
		return
	}

	fmt.Printf("%s\n%s", respNonProf.Name, nonProf.Name)
	assert.Equal(t, respNonProf.Name, nonProf.Name)
	assert.Equal(t, respNonProf.Owner, user.IsOwner)

	type ParamsDonation struct {
		Amount        string `json:"amount"`
		Email         string `json:"email"`
		NonProfitName string `json:"non_profit_name"`
	}
	requestDon, err := sendPostRequest("/api/donations", ParamsDonation{
		Amount:        "400",
		Email:         user.Email,
		NonProfitName: nonProf.Name},
		token)
	if err != nil {
		t.Fatal("Error sending post to create donation")
		return
	}
	assert.Equal(t, 200, requestDon.StatusCode)
}

func sendPostRequest(endpoint string, payload any, token string) (*http.Response, error) {
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	return client.Do(req)
}
func parseJSONResponse(resp *http.Response, target any) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.Unmarshal(body, target)
}
