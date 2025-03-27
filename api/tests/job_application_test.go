package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
	"github.com/stretchr/testify/assert"
)

func TestJobApplications(t *testing.T) {
	queries.Purge(ctx)

	setUpUser(ctx)

	user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

	softwareEngineer, _ := queries.CreateJobApplication(ctx, db.CreateJobApplicationParams{
		UserID:        user.ID,
		CompanyName:   "Evil Corp Inc.",
		JobTitle:      "Software Engineer",
		DateApplied:   pgtype.Timestamptz{Time: time.Now().Add(time.Hour * -24 * 2), Valid: true},
		Status:        db.StatusINPROGRESS,
		MinSalary:     pgtype.Float8{Float64: 50_000.00, Valid: true},
		MaxSalary:     pgtype.Float8{Float64: 70_000.00, Valid: true},
		JobPostingUrl: pgtype.Text{String: "https://glassbore.com/jobs/swe420692137", Valid: true},
	})
	iOSDeveloper, _ := queries.CreateJobApplication(ctx, db.CreateJobApplicationParams{
		UserID:        user.ID,
		CompanyName:   "Apple",
		JobTitle:      "iOS Developer",
		DateApplied:   pgtype.Timestamptz{Time: time.Now().Add(time.Hour * -24), Valid: true},
		Status:        db.StatusREJECTED,
		MinSalary:     pgtype.Float8{Float64: 100_000.00, Valid: true},
		MaxSalary:     pgtype.Float8{Float64: 125_000.00, Valid: true},
		JobPostingUrl: pgtype.Text{String: "https://glassbore.com/jobs/ios420692137", Valid: true},
	})
	angularDeveloper, _ := queries.CreateJobApplication(ctx, db.CreateJobApplicationParams{
		UserID:        user.ID,
		CompanyName:   "Google",
		JobTitle:      "Angular Developer",
		DateApplied:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Status:        db.StatusACCEPTED,
		MinSalary:     pgtype.Float8{Float64: 70_000.00, Valid: true},
		MaxSalary:     pgtype.Float8{Float64: 90_000.00, Valid: true},
		JobPostingUrl: pgtype.Text{String: "https://glassbore.com/jobs/fe420692137", Valid: true},
	})

	t.Run("valid request - default query params", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - second page", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?page=1&size=1", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 1, resBodyRaw.Page)
		assert.Equal(t, 1, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 1)
	})

	t.Run("valid request - third page", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?page=2&size=1", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 2, resBodyRaw.Page)
		assert.Equal(t, 1, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 1)
	})

	t.Run("valid request - sort ascending by company name", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=company_name", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort descending by company name", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=-company_name", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort ascending by job title", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=job_title", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort descending by job title", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=-job_title", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort ascending by date applied", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=date_applied", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort descending by date applied", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=-date_applied", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort ascending by status", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=status", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort descending by status", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=-status", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort ascending by salary", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=salary", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	t.Run("valid request - sort descending by salary", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?sort=-salary", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 3, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 3)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[2].ID, "missing job application id")
		assert.Equal(t, softwareEngineer.ID.String(), resBodyRaw.Data[2].ID)
		assert.Equal(t, softwareEngineer.CompanyName, resBodyRaw.Data[2].CompanyName)
		assert.Equal(t, softwareEngineer.JobTitle, resBodyRaw.Data[2].JobTitle)
		assert.Equal(t, softwareEngineer.DateApplied.Time.UTC(), resBodyRaw.Data[2].DateApplied.UTC())
		assert.Equal(t, softwareEngineer.Status, resBodyRaw.Data[2].Status)
		assert.Equal(t, false, resBodyRaw.Data[2].IsReplied)
		assert.Equal(t, softwareEngineer.MinSalary.Float64, resBodyRaw.Data[2].MinSalary)
		assert.Equal(t, softwareEngineer.MaxSalary.Float64, resBodyRaw.Data[2].MaxSalary)
		assert.Equal(t, softwareEngineer.JobPostingUrl.String, resBodyRaw.Data[2].JobPostingURL)
	})

	// TODO: Test sorting by `is_replied`

	t.Run("valid request - filter by company name", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?company_name_or_job_title=Apple", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 1, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 1)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)
	})

	t.Run("valid request - filter by job title", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?company_name_or_job_title=Dev", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 2, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 2)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, angularDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, angularDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, angularDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, angularDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, angularDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, angularDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, angularDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, angularDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)

		assert.NotEmpty(t, resBodyRaw.Data[1].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[1].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[1].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[1].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[1].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[1].Status)
		assert.Equal(t, false, resBodyRaw.Data[1].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[1].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[1].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[1].JobPostingURL)
	})

	t.Run("valid request - filter by date applied", func(t *testing.T) {
		w := httptest.NewRecorder()

		dateApplied := iOSDeveloper.DateApplied.Time.Format("2006-01-02")

		req, _ := http.NewRequest("GET", "/api/job-applications?date_applied="+dateApplied, nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 1, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 1)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)
	})

	t.Run("valid request - filter by status", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications?status=REJECTED", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationsResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 0, resBodyRaw.Page)
		assert.Equal(t, 10, resBodyRaw.Size)
		assert.Equal(t, 1, resBodyRaw.Total)
		assert.Len(t, resBodyRaw.Data, 1)

		assert.NotEmpty(t, resBodyRaw.Data[0].ID, "missing job application id")
		assert.Equal(t, iOSDeveloper.ID.String(), resBodyRaw.Data[0].ID)
		assert.Equal(t, iOSDeveloper.CompanyName, resBodyRaw.Data[0].CompanyName)
		assert.Equal(t, iOSDeveloper.JobTitle, resBodyRaw.Data[0].JobTitle)
		assert.Equal(t, iOSDeveloper.DateApplied.Time.UTC(), resBodyRaw.Data[0].DateApplied.UTC())
		assert.Equal(t, iOSDeveloper.Status, resBodyRaw.Data[0].Status)
		assert.Equal(t, false, resBodyRaw.Data[0].IsReplied)
		assert.Equal(t, iOSDeveloper.MinSalary.Float64, resBodyRaw.Data[0].MinSalary)
		assert.Equal(t, iOSDeveloper.MaxSalary.Float64, resBodyRaw.Data[0].MaxSalary)
		assert.Equal(t, iOSDeveloper.JobPostingUrl.String, resBodyRaw.Data[0].JobPostingURL)
	})
}

func TestJobApplication(t *testing.T) {
	queries.Purge(ctx)

	setUpUser(ctx)

	user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

	var (
		companyName   = "Evil Corp Inc."
		jobTitle      = "Software Engineer"
		dateApplied   = time.Now().Add(time.Hour * -1)
		status        = db.StatusINPROGRESS
		isReplied     = false
		minSalary     = 50_000.00
		maxSalary     = 70_000.00
		jobPostingURL = "https://glassbore.com/jobs/swe420692137"
		notes         = "Follow up in two weeks"
	)

	jobApplication, _ := queries.CreateJobApplication(ctx, db.CreateJobApplicationParams{
		UserID:        user.ID,
		CompanyName:   companyName,
		JobTitle:      jobTitle,
		DateApplied:   pgtype.Timestamptz{Time: dateApplied, Valid: true},
		Status:        status,
		MinSalary:     pgtype.Float8{Float64: minSalary, Valid: true},
		MaxSalary:     pgtype.Float8{Float64: maxSalary, Valid: true},
		JobPostingUrl: pgtype.Text{String: jobPostingURL, Valid: true},
		Notes:         pgtype.Text{String: notes, Valid: true},
	})

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, status, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("non-existing job application", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/api/job-applications/f4d15edc-e780-42b5-957d-c4352401d9ca", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NotEmpty(t, resBodyRaw.Error, "missing error message")
		assert.Contains(t, resBodyRaw.Error, "no rows in result set")
	})
}

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
			isReplied     = false
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
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
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
		assert.Equal(t, jobApplication.IsReplied, resBodyRaw.IsReplied)
		assert.Equal(t, jobApplication.MinSalary.Float64, resBodyRaw.MinSalary)
		assert.Equal(t, jobApplication.MaxSalary.Float64, resBodyRaw.MaxSalary)
		assert.Equal(t, jobApplication.JobPostingUrl.String, resBodyRaw.JobPostingURL)
		assert.Equal(t, jobApplication.Notes.String, resBodyRaw.Notes)
	})

	t.Run("invalid payload - missing company name", func(t *testing.T) {
		w := httptest.NewRecorder()

		var (
			companyName   = ""
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

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.NotEmpty(t, resBodyRaw.Error, "missing error message")
		assert.Contains(t, resBodyRaw.Error, "CompanyName", "required")
	})

	t.Run("invalid payload - missing job title", func(t *testing.T) {
		w := httptest.NewRecorder()

		var (
			companyName   = "Evil Corp Inc."
			jobTitle      = ""
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

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.NotEmpty(t, resBodyRaw.Error, "missing error message")
		assert.Contains(t, resBodyRaw.Error, "JobTitle", "required")
	})

	t.Run("invalid payload - missing date applied", func(t *testing.T) {
		w := httptest.NewRecorder()

		var (
			companyName   = "Evil Corp Inc."
			jobTitle      = "Software Engineer"
			dateApplied   = time.Time{}
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

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.NotEmpty(t, resBodyRaw.Error, "missing error message")
		assert.Contains(t, resBodyRaw.Error, "DateApplied", "required")
	})

	t.Run("invalid payload - missing status", func(t *testing.T) {
		w := httptest.NewRecorder()

		var (
			companyName             = "Evil Corp Inc."
			jobTitle                = "Software Engineer"
			dateApplied             = time.Now().Add(time.Hour * -1)
			status        db.Status = ""
			minSalary               = 50_000.00
			maxSalary               = 70_000.00
			jobPostingURL           = "https://glassbore.com/jobs/swe420692137"
			notes                   = "Follow up in two weeks"
		)

		bodyRaw := models.NewCreateJobApplicationReqBody(companyName, jobTitle, dateApplied, status, minSalary, maxSalary, jobPostingURL, notes)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("POST", "/api/job-applications", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.NotEmpty(t, resBodyRaw.Error, "missing error message")
		assert.Contains(t, resBodyRaw.Error, "Status", "required")
	})

	t.Run("invalid payload - incorrect status", func(t *testing.T) {
		w := httptest.NewRecorder()

		var (
			companyName             = "Evil Corp Inc."
			jobTitle                = "Software Engineer"
			dateApplied             = time.Now().Add(time.Hour * -1)
			status        db.Status = "UNKNOWN"
			minSalary               = 50_000.00
			maxSalary               = 70_000.00
			jobPostingURL           = "https://glassbore.com/jobs/swe420692137"
			notes                   = "Follow up in two weeks"
		)

		bodyRaw := models.NewCreateJobApplicationReqBody(companyName, jobTitle, dateApplied, status, minSalary, maxSalary, jobPostingURL, notes)
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("POST", "/api/job-applications", strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.Error
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.NotEmpty(t, resBodyRaw.Error, "missing error message")
		assert.Contains(t, resBodyRaw.Error, "Status", "Field validation for 'Status' failed on the 'oneof' tag")
	})
}

func TestUpdateJobApplication(t *testing.T) {
	queries.Purge(ctx)

	setUpUser(ctx)

	user, _ := queries.GetUserByEmail(ctx, "jakub.szewczyk@test.com")

	var (
		companyName   = "Evil Corp Inc."
		jobTitle      = "Software Engineer"
		dateApplied   = time.Now().Add(time.Hour * -1)
		status        = db.StatusINPROGRESS
		isReplied     = false
		minSalary     = 50_000.00
		maxSalary     = 70_000.00
		jobPostingURL = "https://glassbore.com/jobs/swe420692137"
		notes         = "Follow up in two weeks"
	)

	jobApplication, _ := queries.CreateJobApplication(ctx, db.CreateJobApplicationParams{
		UserID:        user.ID,
		CompanyName:   companyName,
		JobTitle:      jobTitle,
		DateApplied:   pgtype.Timestamptz{Time: dateApplied, Valid: true},
		Status:        status,
		MinSalary:     pgtype.Float8{Float64: minSalary, Valid: true},
		MaxSalary:     pgtype.Float8{Float64: maxSalary, Valid: true},
		JobPostingUrl: pgtype.Text{String: jobPostingURL, Valid: true},
		Notes:         pgtype.Text{String: notes, Valid: true},
	})

	t.Run("valid request - changing company name", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			CompanyName: "Google",
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, "Google", resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, status, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing job title", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			JobTitle: "Angular Developer",
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, "Angular Developer", resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, status, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing date applied", func(t *testing.T) {
		w := httptest.NewRecorder()

		dateApplied := time.Date(2006, 02, 01, 0, 0, 0, 0, time.UTC)

		bodyRaw := models.UpdateJobApplicationReqBody{
			DateApplied: dateApplied,
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, status, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing status", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			Status: db.StatusREJECTED,
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, db.StatusREJECTED, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing is replied", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			IsReplied: true,
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, db.StatusREJECTED, resBodyRaw.Status)
		assert.Equal(t, true, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing min salary", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			MinSalary: 2137.00,
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, db.StatusREJECTED, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, 2137.00, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing max salary", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			MaxSalary: 42069.00,
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, db.StatusREJECTED, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, 42069.00, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing job posting URL", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			JobPostingURL: "https://glassbore.com/jobs/fe420692137",
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, db.StatusREJECTED, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, "https://glassbore.com/jobs/fe420692137", resBodyRaw.JobPostingURL)
		assert.Equal(t, notes, resBodyRaw.Notes)
	})

	t.Run("valid request - changing notes", func(t *testing.T) {
		w := httptest.NewRecorder()

		bodyRaw := models.UpdateJobApplicationReqBody{
			Notes: "Follow up in a week",
		}
		bodyJSON, _ := json.Marshal(bodyRaw)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/job-applications/%v", jobApplication.ID), strings.NewReader(string(bodyJSON)))
		req.Header.Add("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		var resBodyRaw models.JobApplicationResBody
		err := json.Unmarshal(w.Body.Bytes(), &resBodyRaw)

		assert.NoError(t, err, "error unmarshaling response body")

		assert.Equal(t, http.StatusOK, w.Code)

		assert.NotEmpty(t, resBodyRaw.ID, "missing job application id")
		assert.Equal(t, companyName, resBodyRaw.CompanyName)
		assert.Equal(t, jobTitle, resBodyRaw.JobTitle)
		assert.Equal(t, dateApplied.UTC(), resBodyRaw.DateApplied.UTC())
		assert.Equal(t, db.StatusREJECTED, resBodyRaw.Status)
		assert.Equal(t, isReplied, resBodyRaw.IsReplied)
		assert.Equal(t, minSalary, resBodyRaw.MinSalary)
		assert.Equal(t, maxSalary, resBodyRaw.MaxSalary)
		assert.Equal(t, jobPostingURL, resBodyRaw.JobPostingURL)
		assert.Equal(t, "Follow up in a week", resBodyRaw.Notes)
	})
}
