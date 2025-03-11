package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	queries.Purge(ctx)

	setupUser(ctx)

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/profile", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ProfileResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing user id")
		assert.Equal(t, "Jakub", resBodyRaw.FirstName)
		assert.Equal(t, "Szewczyk", resBodyRaw.LastName)
		assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw.Email)
		assert.Equal(t, false, resBodyRaw.IsEmailVerified)
	})

	t.Run("missing authorization token", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/profile", nil)

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		assert.Equal(t, "missing Authorization header", resBodyRaw.Error)
	})

	t.Run("invalid authorization token", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/profile", nil)
		req.Header.Add("Authorization", "testing")

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		assert.Equal(t, "invalid Authorization header format", resBodyRaw.Error)
	})

	t.Run("expired authorization token", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/profile", nil)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg5MTQ4NTcsInN1YiI6ImphY2suZGFuaWVsc0BnbWFpbC5jb20iLCJ1aWQiOiI1YWFhMTAzMS01MDZjLTRmYzItYTMzNC1lYTVhOTQzNmYzYmQifQ.xcKER7PtNpujouNS_VlWePIDDHQvAkdO40XckPGcmcs")

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		assert.Equal(t, "token has invalid claims: token is expired", resBodyRaw.Error)
	})
}

func TestVerifyEmail(t *testing.T) {
	queries.Purge(ctx)

	setupUser(ctx)

	user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewVerifyEmailReqBody(user.VerificationToken)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PATCH", "/api/profile/verify-email", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ProfileResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing user id")
		assert.Equal(t, "Jakub", resBodyRaw.FirstName)
		assert.Equal(t, "Szewczyk", resBodyRaw.LastName)
		assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw.Email)
		assert.Equal(t, true, resBodyRaw.IsEmailVerified)
	})

	t.Run("missing verification token", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewVerifyEmailReqBody("")
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PATCH", "/api/profile/verify-email", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Equal(t, "Key: 'VerifyEmailReqBody.VerificationToken' Error:Field validation for 'VerificationToken' failed on the 'required' tag", resBodyRaw.Error)
	})

	t.Run("invalid verification token", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewVerifyEmailReqBody(user.Email)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PATCH", "/api/profile/verify-email", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Equal(t, "invalid verification token", resBodyRaw.Error)
	})

	// TODO: Mock time
	// t.Run("expired verification token", func(t *testing.T) {
	// 	w := httptest.NewRecorder()
	//
	// 	bodyRaw := models.NewVerifyEmailReqBody(user.VerificationToken)
	// 	bodyJSON, _ := json.Marshal(bodyRaw)
	//
	// 	req, _ := http.NewRequest("PATCH", "/api/profile/verify-email", strings.NewReader(string(bodyJSON)))
	// 	req.Header.Add("Authorization", "Bearer "+token)
	//
	// 	r.ServeHTTP(w, req)
	//
	// 	var resBodyRaw models.Error
	// 	err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)
	//
	// 	assert.NoError(t, err, "error unmarshaling response body")
	//
	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	//
	// 	assert.Equal(t, "expired verification token", resBodyRaw.Error)
	// })
}

func TestSendVerificationEmail(t *testing.T) {
	queries.Purge(ctx)

	setupUser(ctx)

	user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		tkn, _ := queries.GetVerificationToken(ctx, user.ID)

		req, _ := http.NewRequest("GET", "/api/profile/verify-email", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		newTkn, _ := queries.GetVerificationToken(ctx, user.ID)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, tkn.Token, newTkn.Token)
		assert.Equal(t, tkn.ExpiresAt, newTkn.ExpiresAt)
	})

	t.Run("expired verification token", func(t *testing.T) {
		w := httptest.NewRecorder()

		queries.ExpireVerificationToken(ctx, user.ID)
		expiredToken, _ := queries.GetVerificationToken(ctx, user.ID)

		req, _ := http.NewRequest("GET", "/api/profile/verify-email", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		renewedToken, _ := queries.GetVerificationToken(ctx, user.ID)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEqual(t, expiredToken.Token, renewedToken.Token)
		assert.NotEqual(t, expiredToken.ExpiresAt.Time.String(), renewedToken.ExpiresAt.Time.String())
	})
}
