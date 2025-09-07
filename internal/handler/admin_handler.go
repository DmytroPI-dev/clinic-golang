package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"

)

func ShowDashboard (c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{"title": "Admin Dashboard"})
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login page"})
}