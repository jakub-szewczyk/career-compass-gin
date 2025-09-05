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

func TestCreateResume(t *testing.T) {
	t.Run("inavlid request - unauthorized", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/resumes", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("valid request - create resume", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		w := httptest.NewRecorder()

		title := "Evil Corp Inc. personalized"

		bodyRaw := models.NewCreateResumeReqBody(title)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("POST", "/api/resumes", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.CreateResumeResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusCreated, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing resume id")
		assert.Equal(t, title, resBodyRaw.Title)
	})

	t.Run("valid request - default resume title", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		w := httptest.NewRecorder()

		bodyRaw := models.NewCreateResumeReqBody("")
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("POST", "/api/resumes", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.CreateResumeResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusCreated, w.Code)

		assert.Contains(t, resBodyRaw.Title, "Untitled 1")
	})

	t.Run("valid request - default resume title increments", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		w1 := httptest.NewRecorder()

		bodyRaw1 := models.NewCreateResumeReqBody("")
		bodyJSON1, _ := json.Marshal(bodyRaw1)

		req1, _ := http.NewRequest("POST", "/api/resumes", strings.NewReader(string(bodyJSON1)))
		req1.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(w1, req1)

		var resBodyRaw1 models.CreateResumeResBody
		err := json.Unmarshal(w1.Body.Bytes(), &resBodyRaw1)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusCreated, w1.Code)

		assert.Contains(t, resBodyRaw1.Title, "Untitled 1")

		w2 := httptest.NewRecorder()

		bodyRaw2 := models.NewCreateResumeReqBody("")
		bodyJSON2, _ := json.Marshal(bodyRaw2)

		req2, _ := http.NewRequest("POST", "/api/resumes", strings.NewReader(string(bodyJSON2)))
		req2.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w2, req2)

		var resBodyRaw2 models.CreateResumeResBody
		err = json.Unmarshal(w2.Body.Bytes(), &resBodyRaw2)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusCreated, w2.Code)

		assert.Contains(t, resBodyRaw2.Title, "Untitled 2")
	})

	t.Run("invalid request - duplicate title", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		title := "Evil Corp Inc. personalized"

		w1 := httptest.NewRecorder()

		bodyRaw1 := models.NewCreateResumeReqBody(title)
		bodyJSON1, _ := json.Marshal(bodyRaw1)

		req1, _ := http.NewRequest("POST", "/api/resumes", strings.NewReader(string(bodyJSON1)))
		req1.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(w1, req1)

		assert.Equal(t, http.StatusCreated, w1.Code)

		w2 := httptest.NewRecorder()

		bodyRaw2 := models.NewCreateResumeReqBody(title)
		bodyJSON2, _ := json.Marshal(bodyRaw2)

		req2, _ := http.NewRequest("POST", "/api/resumes", strings.NewReader(string(bodyJSON2)))
		req2.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusInternalServerError, w2.Code)
	})

	t.Run("invalid request - malformed payload", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/resumes", strings.NewReader("{invalid json}"))
		req.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
