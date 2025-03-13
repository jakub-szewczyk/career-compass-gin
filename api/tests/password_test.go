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

func TestInitPasswordReset(t *testing.T) {
	queries.Purge(ctx)

	setUpUser(ctx)

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewInitPasswordResetReqBody("jakub.szewczyk@test.com")
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("POST", "/api/password/reset", strings.NewReader(string(bodyJSON)))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid payload - missing email", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/password/reset", nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("non-existing user", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewInitPasswordResetReqBody("john.doe@test.com")
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("POST", "/api/password/reset", strings.NewReader(string(bodyJSON)))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestResetPassword(t *testing.T) {
	queries.Purge(ctx)

	user, _ := setUpUser(ctx)

	token, _ := queries.CreatePasswordResetToken(ctx, user.ID)

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewResetPasswordReqBody("CareerCompass!123", "CareerCompass!123", token)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", "/api/password/reset", strings.NewReader(string(bodyJSON)))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		_, err := queries.GetPasswordResetToken(ctx, token)

		assert.Error(t, err)
	})

	t.Run("invalid payload - mismatching passwords", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewResetPasswordReqBody("CareerCompass!123", "qwerty!123456789", token)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", "/api/password/reset", strings.NewReader(string(bodyJSON)))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid payload - missing password reset token", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.NewResetPasswordReqBody("CareerCompass!123", "CareerCompass!123", "")
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", "/api/password/reset", strings.NewReader(string(bodyJSON)))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("non-existing password reset token", func(t *testing.T) {
		queries.DeletePasswordResetToken(ctx, token)

		w := httptest.NewRecorder()

		bodyRaw := models.NewResetPasswordReqBody("CareerCompass!123", "CareerCompass!123", token)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", "/api/password/reset", strings.NewReader(string(bodyJSON)))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
