package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// Render new program template
func AdminShowNewProgramForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "program-form.html", gin.H{
		"Categories": models.AllCategories,
	})
}

// Create new program template
func AdminCreateNewProgram(db *gorm.DB) gin.HandlerFunc {
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
			log.Printf("Failed to create program: %s", err)
			ctx.Status(http.StatusInternalServerError)
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

// ShowEditProgramForm finds a program by ID and renders the edit form.
func AdminShowEditProgramForm(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Find the program in the database
		var program models.Program
		if err := db.First(&program, id).Error; err != nil {
			// Handle the case where no record found
			log.Printf("Failed to find program with ID %s: %s", id, err)
			ctx.Status(http.StatusNotFound)
			return
		}
		// Render the edit form with the program data
		ctx.HTML(http.StatusOK, "program-form.html", gin.H{
			"Categories": models.AllCategories,
			"Program":    program,
		})
	}
}

// UpdateProgram handles the submission of the edit program form.
func AdminUpdateProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")

		// Find existing program
		var program models.Program
		if err := db.First(&program, id).Error; err != nil {
			// Handle the case where no record found
			log.Printf("Failed to find program with ID %s: %s", id, err)
			ctx.Status(http.StatusNotFound)
			return
		}

		//Parse data from the request
		program.Title = ctx.PostForm("title")
		program.Description = ctx.PostForm("description")
		program.Results = ctx.PostForm("results")
		program.Category = ctx.PostForm("category")
		// Will update translation fields later

		// Save updates to the DB
		if err := db.Save(&program).Error; err != nil {
			log.Printf("Failed to update program with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return the updated program
		ctx.HTML(http.StatusOK, "program-row.html", program)
	}
}


