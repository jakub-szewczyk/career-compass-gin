package handlers

import (
	"net/http"
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
}

type signInReqBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=16"` // TODO: Improve password strength
}

type signInResBody struct {
	User  db.CreateUserRow `json:"user"`
	Token string           `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var body signInReqBody

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

	c.JSON(http.StatusCreated, signInResBody{
		User: db.CreateUserRow{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Token: signed,
	})
}
