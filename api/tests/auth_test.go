package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	queries.Purge(ctx)

	w := httptest.NewRecorder()

	bodyRaw := models.NewSignUpReqBody("Jakub", "Szewczyk", "jakub.szewczyk@test.com", "qwerty!123456789", "qwerty!123456789")
	bodyJSON, _ := json.Marshal(bodyRaw)

	req, _ := http.NewRequest("POST", "/api/sign-up", strings.NewReader(string(bodyJSON)))

	r.ServeHTTP(w, req)

	var resBodyRaw models.SignUpResBody
	err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

	token = resBodyRaw.Token

	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.NotEmpty(t, resBodyRaw.User.ID, "missing user id")
	assert.Equal(t, "Jakub", resBodyRaw.User.FirstName)
	assert.Equal(t, "Szewczyk", resBodyRaw.User.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw.User.Email)
	assert.Equal(t, false, resBodyRaw.User.IsEmailVerified)
	assert.NotEmpty(t, resBodyRaw.Token, "missing token")

	uuid, _ := common.ToUUID(resBodyRaw.User.ID)
	user, err := queries.GetUserById(ctx, uuid)

	assert.NoError(t, err, "error getting user from the database")

	assert.NotEmpty(t, user.ID, "missing user id")
	assert.Equal(t, "Jakub", user.FirstName)
	assert.Equal(t, "Szewczyk", user.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", user.Email)
	assert.Equal(t, false, user.IsEmailVerified.Bool)
}

func TestSignIn(t *testing.T) {
	queries.Purge(ctx)

	TestSignUp(t)

	w := httptest.NewRecorder()

	bodyRaw := models.NewSignInReqBody("jakub.szewczyk@test.com", "qwerty!123456789")
	bodyJSON, _ := json.Marshal(bodyRaw)

	req, _ := http.NewRequest("POST", "/api/sign-in", strings.NewReader(string(bodyJSON)))

	r.ServeHTTP(w, req)

	var resBodyRaw models.SignInResBody
	err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotEmpty(t, resBodyRaw.User.ID, "missing user id")
	assert.Equal(t, "Jakub", resBodyRaw.User.FirstName)
	assert.Equal(t, "Szewczyk", resBodyRaw.User.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw.User.Email)
	assert.Equal(t, false, resBodyRaw.User.IsEmailVerified)
	assert.NotEmpty(t, resBodyRaw.Token, "missing token")
}
