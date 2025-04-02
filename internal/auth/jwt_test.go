package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const TOKEN_SECRET_TEST = "mylittlesecret"

func TestJWT(t *testing.T) {
	userIDtest := uuid.New()
	token, err := MakeJWT(userIDtest, TOKEN_SECRET_TEST, time.Duration(time.Hour))
	if err != nil {
		t.Fatal("Can't create token")
	}

	id, err := ValidateJWT(token, TOKEN_SECRET_TEST)
	if err != nil {
		t.Fatal("Token not validated")
	}
	assert.Equal(t, userIDtest, id)
}
