package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"golang.org/x/crypto/bcrypt"
)

type signUpReqBody struct {
	FirstName       string `json:"firstName" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=16"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

type signUpResBody struct {
	Token string           `json:"token"`
	User  db.CreateUserRow `json:"user"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var body signUpReqBody

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

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.ID,
		"sub": user.Email,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	k := os.Getenv("JWT_SECRET")

	s, err := t.SignedString([]byte(k))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, signUpResBody{
		Token: s,
		User:  user,
	})
}
