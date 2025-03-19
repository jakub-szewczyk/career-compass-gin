package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
	"github.com/stretchr/testify/assert"
)

// TODO: Test getting job application by id

func TestCreateJobApplication(t *testing.T) {
	queries.Purge(ctx)

	setUpUser(ctx)

	user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		var (
			companyName   = "Evil Corp Inc."
			jobTitle      = "Software Engineer"
			dateApplied   = time.Now().Add(time.Hour * -1)
			status        = db.StatusINPROGRESS
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
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, status, resBodyRaw.Status)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)

		uuid, err := common.ToUUID(resBodyRaw.ID)

		jobApplication, err := queries.GetJobApplication(ctx, db.GetJobApplicationParams{
			ID:     uuid,
			UserID: user.ID,
		})

		assert.NoError(t, err, "error querying job application")
		assert.NotEmpty(t, jobApplication.ID, "missing job application id")
		assert.Equal(t, jobApplication.CompanyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobApplication.JobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, jobApplication.DateApplied.Time.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, jobApplication.Status, resBodyRaw.Status)
		assert.Equal(t, jobApplication.MinSalary.Float64, resBodyRaw.MinSalary)
		assert.Equal(t, jobApplication.MaxSalary.Float64, resBodyRaw.MaxSalary)
		assert.Equal(t, jobApplication.JobPostingUrl.String, resBodyRaw.JobPostingURL)
		assert.Equal(t, jobApplication.Notes.String, resBodyRaw.Notes)
	})

	// TODO: Missing suites
}
