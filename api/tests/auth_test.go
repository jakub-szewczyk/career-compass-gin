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

	R.ServeHTTP(w, req)

	var resBodyRaw handlers.SignUpResBody
	err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.NoError(t, err, "error unmarshaling response body")

	assert.NotEmpty(t, resBodyRaw.User.ID, "missing user id")

	assert.Equal(t, "Jakub", resBodyRaw.User.FirstName)
	assert.Equal(t, "Szewczyk", resBodyRaw.User.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw.User.Email)

	assert.NotEmpty(t, resBodyRaw.Token, "missing token")

	// TODO: Test database for new entry
}
