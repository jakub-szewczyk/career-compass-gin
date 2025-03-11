package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
)

// Profile godoc
//
//	@Summary		Get user profile
//	@Description	Retrieves and returns the profile information of the currently authenticated user
//
//	@Security		BearerAuth
//
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Success		200	{object}	models.ProfileResBody
//	@Router			/profile [get]
func (h *Handler) Profile(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.queries.GetUserById(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	profileResBody, err := models.NewProfileResBody(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profileResBody)
}

// SendVerificationEmail godoc
//
//	@Summary		Send user verification email
//	@Description	Sends a verification email to the user. This endpoint can be used to resend the email if needed.
//
//	@Security		BearerAuth
//
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Success		204
//	@Router			/profile/verify-email [get]
func (h *Handler) SendVerificationEmail(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.queries.GetVerificationToken(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if token.ExpiresAt.Time.Before(time.Now()) {
		newToken, err := h.queries.UpdateVerificationToken(h.ctx, uuid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		token.Token = newToken.Token
		token.ExpiresAt = newToken.ExpiresAt
	}

	user, err := h.queries.GetUserById(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	}

	// TODO: Consider using goroutines
	tmpl, err := template.ParseFiles(filepath.Join("templates", "sign-up.html"))
	if err != nil {
		fmt.Println("error loading template:", err)
		return
	}

	var html bytes.Buffer
	if err := tmpl.Execute(&html, struct {
		FirstName string
		Link      string
		Year      int
	}{
		FirstName: user.FirstName,
		Link:      h.env.EmailVerificationURL + fmt.Sprintf("?token=%v", token.Token),
		Year:      time.Now().Year(),
	}); err != nil {
		fmt.Println("error rendering template:", err)
		return
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: Welcome to Career Compass, %v!\n", user.FirstName)

	auth := smtp.PlainAuth(h.env.SMTPIdentity, h.env.SMTPUsername, h.env.SMTPPassword, h.env.SMTPHost)
	err = smtp.SendMail(h.env.SMTPHost+":"+h.env.SMTPPORT, auth, h.env.SMTPUsername, []string{user.Email}, []byte(subject+mime+html.String()))
	if err != nil {
		fmt.Println("error sending email:", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// VeifyEmail godoc
//
//	@Summary		Verify user email
//	@Description	Confirms the email address of the currently authenticated user. This endpoint requires an email verification token sent to the user's registered email.
//
//	@Security		BearerAuth
//
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.VerifyEmailReqBody	true	"Email verification data"
//	@Failure		400		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		200		{object}	models.ProfileResBody
//	@Router			/profile/verify-email [patch]
func (h *Handler) VerifyEmail(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body models.VerifyEmailReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.queries.GetVerificationToken(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if body.VerificationToken != token.Token {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid verification token",
		})
		return
	}

	if token.ExpiresAt.Time.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "expired verification token",
		})
		return
	}

	user, err := h.queries.VerifyEmail(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	profileResBody, err := models.NewProfileResBody(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profileResBody)
}
