package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/stretchr/testify/assert"
)

func TestResumes(t *testing.T) {
	queries.Purge(ctx)

	setUpUser(ctx)

	user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

	resume1, _ := queries.CreateResume(ctx, db.CreateResumeParams{
		UserID: user.ID,
		Title:  "Resume A",
	})
	resume2, _ := queries.CreateResume(ctx, db.CreateResumeParams{
		UserID: user.ID,
		Title:  "Resume C",
	})
	resume3, _ := queries.CreateResume(ctx, db.CreateResumeParams{
		UserID: user.ID,
		Title:  "Resume B",
	})

	t.Run("valid request - default query params", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)

		assert.Len(t, resBodyRaw.Data, 3)

		assert.Equal(t, resume3.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[2].ID)
	})

	t.Run("valid request - pagination", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?page=1&size=1", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 1, resBodyRaw.Page)
		assert.Equal(t, 1, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)

		assert.Len(t, resBodyRaw.Data, 1)

		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[0].ID)
	})

	t.Run("valid request - filter by title", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?title=Resume A", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 1, resBodyRaw.Total)

		assert.Len(t, resBodyRaw.Data, 1)

		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[0].ID)
	})

	t.Run("valid request - sort ascending by title", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?sort=title", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Len(t, resBodyRaw.Data, 3)

		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, resume3.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[2].ID)
	})

	t.Run("valid request - sort descending by title", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?sort=-title", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Len(t, resBodyRaw.Data, 3)

		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, resume3.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[2].ID)
	})

	t.Run("valid request - sort ascending by date created", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?sort=created_at", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Len(t, resBodyRaw.Data, 3)

		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, resume3.ID.String(), resBodyRaw.Data[2].ID)
	})

	t.Run("valid request - sort descending by date created", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?sort=-created_at", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Len(t, resBodyRaw.Data, 3)

		assert.Equal(t, resume3.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[2].ID)
	})

	t.Run("valid request - sort ascending by last modified", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?sort=updated_at", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Len(t, resBodyRaw.Data, 3)

		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, resume3.ID.String(), resBodyRaw.Data[2].ID)
	})

	t.Run("valid request - sort descending by last modified", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/resumes?sort=-updated_at", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.ResumesResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Len(t, resBodyRaw.Data, 3)

		assert.Equal(t, resume3.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, resume2.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, resume1.ID.String(), resBodyRaw.Data[2].ID)
	})
}

func TestCreateResume(t *testing.T) {
	t.Run("invalid request - unauthorized", func(t *testing.T) {
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

func TestDeleteResume(t *testing.T) {
	const title = "Evil Corp Inc. personalized"

	t.Run("valid request - delete resume", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

		resume, err := queries.CreateResume(ctx, db.CreateResumeParams{UserID: user.ID, Title: title})

		assert.NoError(t, err, "error creating resume")

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", "/api/resumes/"+resume.ID.String(), nil)
		req.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)

		var resBodyRaw models.DeleteResumeResBody
		err = json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing resume id")
		assert.Equal(t, title, resBodyRaw.Title)

		// TODO: Uncomment once `queries.GetResumes` is implemented
		// resumes, err := queries.GetResumes(ctx, db.GetResumesParams{})
		// assert.Len(t, resumes, 0)
	})

	t.Run("invalid request - unauthorized", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

		resume, err := queries.CreateResume(ctx, db.CreateResumeParams{UserID: user.ID, Title: title})

		assert.NoError(t, err)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", "/api/resumes/"+resume.ID.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request - non-existent resume", func(t *testing.T) {
		queries.Purge(ctx)

		setUpUser(ctx)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", "/api/resumes/3912beb6-cb36-4190-a543-8ab4a3f29d4d", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
