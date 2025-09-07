package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ShowNewProgramForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "program-form.html", gin.H{
		"Categories": models.AllCategories,
	})
}

func CreateNewProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse form data from the request
		title := ctx.PostForm("title")
		description := ctx.PostForm("description")
		results := ctx.PostForm("results")
		category := ctx.PostForm("category")
		// Create a new program model instance with the data
		newProgram := models.Program{
			Title:         title,
			Description:   description,
			Results:       results,
			Category:      category,
			TitlePL:       title,
			TitleEN:       title,
			TitleUK:       title,
			DescriptionPL: description,
			DescriptionEN: description,
			DescriptionUK: description,
			ResultsPL:     results,
			ResultsEN:     results,
			ResultsUK:     results,
		}
		// Save the newly created program to DB
		if err := db.Create(&newProgram).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create program"})
			return
		}
		// REnder and return HTML fragment for new row
		ctx.HTML(http.StatusOK, "program-row.html", newProgram)
	}
}
