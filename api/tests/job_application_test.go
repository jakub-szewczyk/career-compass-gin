package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateJobApplication(t *testing.T) {
	queries.Purge(ctx)

	setUpUser(ctx)

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		var (
			companyName   = "Evil Corp Inc."
			jobTitle      = "Software Engineer"
			dateApplied   = time.Now().Add(time.Hour * -1)
			status        = models.IN_PROGRESS
			minSalary     = 50_000.00
			maxSalary     = 70_000.00
			jobPostingURL = "https://glassbore.com/jobs/swe420692137"
			notes         = "Follow up in two weeks"
		)

		bodyRaw := models.NewCreateJobApplicationReqBody(companyName, jobTitle, dateApplied, status, minSalary, maxSalary, jobPostingURL, notes)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("POST", "/api/job-applications", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.CreateJobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusCreated, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied, resBodyRaw.DateApplied)
		assert.Equal(t, status, resBodyRaw.Status)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)

		// TODO: Inspect database
	})
}
