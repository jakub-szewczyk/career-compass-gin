package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jakub-szewczyk/career-compass-gin/api/handlers"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	queries.Purge(ctx)

	w := httptest.NewRecorder()

	bodyRaw := handlers.SignUpReqBody{
		FirstName:       "Jakub",
		LastName:        "Szewczyk",
		Email:           "jakub.szewczyk@test.com",
		Password:        "qwerty!123456789",
		ConfirmPassword: "qwerty!123456789",
	}
	bodyJSON, _ := json.Marshal(bodyRaw)

	req, _ := http.NewRequest("POST", "/api/sign-up", strings.NewReader(string(bodyJSON)))

	r.ServeHTTP(w, req)

	var resBodyRaw handlers.SignUpResBody
	err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

	token = resBodyRaw.Token

	// NOTE: Test response body
	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.NotEmpty(t, resBodyRaw.User.ID, "missing user id")
	assert.Equal(t, "Jakub", resBodyRaw.User.FirstName)
	assert.Equal(t, "Szewczyk", resBodyRaw.User.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw.User.Email)
	assert.Equal(t, false, resBodyRaw.User.IsEmailVerified.Bool)
	assert.NotEmpty(t, resBodyRaw.User.VerificationToken, "missing email verification token")
	assert.NotEmpty(t, resBodyRaw.Token, "missing token")

	// NOTE: Test database entry
	user, err := queries.GetUserById(ctx, resBodyRaw.User.ID)

	assert.NoError(t, err, "error getting user from the database")

	assert.NotEmpty(t, user.ID, "missing user id")
	assert.Equal(t, "Jakub", user.FirstName)
	assert.Equal(t, "Szewczyk", user.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", user.Email)
	assert.Equal(t, false, user.IsEmailVerified.Bool)
	assert.NotEmpty(t, user.VerificationToken, "missing email verification token")
}

func TestSignIn(t *testing.T) {
	queries.Purge(ctx)

	TestSignUp(t)

	w := httptest.NewRecorder()

	bodyRaw := handlers.SignInReqBody{
		Email:    "jakub.szewczyk@test.com",
		Password: "qwerty!123456789",
	}
	bodyJSON, _ := json.Marshal(bodyRaw)

	req, _ := http.NewRequest("POST", "/api/sign-in", strings.NewReader(string(bodyJSON)))

	r.ServeHTTP(w, req)

	var resBodyRaw handlers.SignInResBody
	err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

	// NOTE: Test response body
	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotEmpty(t, resBodyRaw.User.ID, "missing user id")
	assert.Equal(t, "Jakub", resBodyRaw.User.FirstName)
	assert.Equal(t, "Szewczyk", resBodyRaw.User.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw.User.Email)
	assert.Equal(t, false, resBodyRaw.User.IsEmailVerified.Bool)
	assert.NotEmpty(t, resBodyRaw.User.VerificationToken, "missing email verification token")
	assert.NotEmpty(t, resBodyRaw.Token, "missing token")
}
