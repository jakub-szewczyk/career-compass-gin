package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/db"
)

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.queries.GetUsers(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if users == nil {
		users = []db.GetUsersRow{}
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
