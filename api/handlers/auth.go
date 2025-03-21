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
	"github.com/golang-jwt/jwt/v5"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"golang.org/x/crypto/bcrypt"
)

// SignUp godoc
//
//	@Summary		User sign up
//	@Description	Registers a new user account with the provided details, including email, password, and other relevant information. Verification email will be sent.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.SignUpReqBody	true	"User sign up data"
//	@Failure		400		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		201		{object}	models.SignUpResBody
//	@Router			/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var body models.SignUpReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.queries.CreateUser(h.ctx, db.CreateUserParams{
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Password:  string(hash),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.ID,
		"sub": user.Email,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	signed, err := token.SignedString([]byte(h.env.JWTSecret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resBody, err := models.NewSignUpResBody(user, signed)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resBody)

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
		Link:      h.env.EmailVerificationURL + fmt.Sprintf("?token=%v", user.VerificationToken),
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
}

// SignIn godoc
//
//	@Summary		User sign in
//	@Description	Authenticates a user and returns a JWT token for session management. Valid credentials are required to access the system.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.SignInReqBody	true	"User sign in data"
//	@Failure		400		{object}	models.Error
//	@Failure		401		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		200		{object}	models.SignInResBody
//	@Router			/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var body models.SignInReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.queries.GetUserOnSignIn(h.ctx, body.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.ID,
		"sub": user.Email,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	signed, err := token.SignedString([]byte(h.env.JWTSecret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resBody, err := models.NewSignInResBody(user, signed)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resBody)
}
