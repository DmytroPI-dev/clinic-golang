package handler

import (
	"log"
	"net/http"

	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Render new program template
func ShowNewProgramForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "program-form.html", gin.H{
		"Categories": models.AllCategories,
	})
}

// Create new program template
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

// AdminDeleteProgram handles the deletion of a program from the admin panel.
func AdminDeleteProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Delete the program from the database
		if err := db.Delete(&models.Program{}, id).Error; err != nil {
			// We will just log the error now, later adding flash error
			log.Printf("Failed to delete program with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return an empty response
		ctx.String(http.StatusOK, "")
	}
}
