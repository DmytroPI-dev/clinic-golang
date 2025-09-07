package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowNewProgramForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "program-form.html", gin.H{
		"Categories": models.AllCategories,
	})
}
