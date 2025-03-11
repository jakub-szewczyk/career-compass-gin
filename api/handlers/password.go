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
)

// InitPasswordReset godoc
//
//	@Summary		Initiate password reset
//	@Description	Generates and sends a password reset token to the user's email address
//	@Tags			Password
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.InitPasswordResetReqBody	true	"User's email address"
//	@Failure		400		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		204
//	@Router			/password/reset [post]
func (h *Handler) InitPasswordReset(c *gin.Context) {
	var body models.InitPasswordResetReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.queries.GetUserByEmail(h.ctx, body.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.queries.CreatePasswordResetToken(h.ctx, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)

	// TODO: Consider using goroutines
	tmpl, err := template.ParseFiles(filepath.Join("templates", "reset-password.html"))
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
		Link:      h.env.ResetPasswordURL + fmt.Sprintf("?token=%v", token),
		Year:      time.Now().Year(),
	}); err != nil {
		fmt.Println("error rendering template:", err)
		return
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Reset Your Password\n"

	auth := smtp.PlainAuth(h.env.SMTPIdentity, h.env.SMTPUsername, h.env.SMTPPassword, h.env.SMTPHost)
	err = smtp.SendMail(h.env.SMTPHost+":"+h.env.SMTPPORT, auth, h.env.SMTPUsername, []string{user.Email}, []byte(subject+mime+html.String()))
	if err != nil {
		fmt.Println("error sending email:", err)
		return
	}
}
