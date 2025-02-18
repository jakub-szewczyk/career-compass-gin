package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	queries.Purge(ctx)

	TestSignUp(t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/profile", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	var resBodyRaw1 db.GetUserByIdRow
	err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw1)

	// NOTE: Test response body
	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotEmpty(t, resBodyRaw1.ID, "missing user id")
	assert.Equal(t, "Jakub", resBodyRaw1.FirstName)
	assert.Equal(t, "Szewczyk", resBodyRaw1.LastName)
	assert.Equal(t, "jakub.szewczyk@test.com", resBodyRaw1.Email)
	assert.Equal(t, false, resBodyRaw1.IsEmailVerified.Bool)

	// NOTE: Test missing Authorization token
	w = httptest.NewRecorder()

	req.Header.Del("Authorization")

	r.ServeHTTP(w, req)

	var resBodyRaw2 struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resBodyRaw2)

	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	assert.Equal(t, "missing Authorization header", resBodyRaw2.Error)

	// NOTE: Test invalid Authorization token
	w = httptest.NewRecorder()

	req.Header.Set("Authorization", "testing")

	r.ServeHTTP(w, req)

	var resBodyRaw3 struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resBodyRaw3)

	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	assert.Equal(t, "invalid Authorization header format", resBodyRaw3.Error)

	// NOTE: Test expired Authorization token
	w = httptest.NewRecorder()

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg5MTQ4NTcsInN1YiI6ImphY2suZGFuaWVsc0BnbWFpbC5jb20iLCJ1aWQiOiI1YWFhMTAzMS01MDZjLTRmYzItYTMzNC1lYTVhOTQzNmYzYmQifQ.xcKER7PtNpujouNS_VlWePIDDHQvAkdO40XckPGcmcs")

	r.ServeHTTP(w, req)

	var resBodyRaw4 struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resBodyRaw4)

	assert.NoError(t, err, "error unmarshaling response body")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	assert.Equal(t, "token has invalid claims: token is expired", resBodyRaw4.Error)
}
