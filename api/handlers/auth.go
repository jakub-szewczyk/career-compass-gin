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
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"golang.org/x/crypto/bcrypt"
)

type SignUpReqBody struct {
	FirstName       string `json:"firstName" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=16"` // TODO: Improve password strength
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

type SignUpResBody struct {
	User  db.CreateUserRow `json:"user"`
	Token string           `json:"token"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var body SignUpReqBody

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

	c.JSON(http.StatusCreated, SignUpResBody{
		User:  user,
		Token: signed,
	})

	tmpl, err := template.ParseFiles(filepath.Join("email", "sign-up.html"))
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
		Link:      h.env.FrontendURL,
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

type SignInReqBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=16"` // TODO: Improve password strength
}

type SignInResBody struct {
	User  db.CreateUserRow `json:"user"`
	Token string           `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var body SignInReqBody

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

	c.JSON(http.StatusOK, SignInResBody{
		User: db.CreateUserRow{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Token: signed,
	})
}
