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

	setupUser(ctx)

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
