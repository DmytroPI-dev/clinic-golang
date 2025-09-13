package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Render new program template
func AdminShowNewProgramForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "program-form.html", gin.H{
		"Categories": models.AllCategories,
		"Program":    models.Program{},
	})
}

// Rendering programs
func ShowProgramsPage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var programs []models.Program
		db.Order("id asc").Find(&programs)
		session := sessions.Default(c)
		userName := session.Get("userName")
		userRole := session.Get("userRole")

		// Render the specific page template. It will handle the layout.
		c.HTML(http.StatusOK, "programs.html", gin.H{
			"Title":    "Manage Programs",
			"User":     userName,
			"UserRole": userRole,
			"Items":    programs,
		})
	}
}

// Create new program template
func AdminCreateNewProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newProgram models.Program
		if err := ctx.ShouldBind(&newProgram); err != nil {
			log.Printf("Failed to bind program data: %s", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		// Save the newly created program to DB
		if err := db.Create(&newProgram).Error; err != nil {
			log.Printf("Failed to create program: %s", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Get user role from session to correctly render the row template
		session := sessions.Default(ctx)
		userRole := session.Get("userRole")
		// Render and return HTML fragment for new row, passing data in a map
		ctx.HTML(http.StatusOK, "program-row.html", gin.H{
			"Item":     newProgram,
			"UserRole": userRole,
		})
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.HTML(http.StatusNotFound, "404.html", gin.H{"Title": "Not Found"})
			} else {
				log.Printf("Failed to find program with ID %s: %s", id, err)
				ctx.Status(http.StatusNotFound)
			}
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
		// Bind form data to the existing program struct
		if err := ctx.ShouldBind(&program); err != nil {
			log.Printf("Failed to bind program data: %s", err)
			ctx.Status(http.StatusBadRequest)
			return
		}

		// Save updates to the DB
		if err := db.Save(&program).Error; err != nil {
			log.Printf("Failed to update program with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Get user role from session to correctly render the row template
		session := sessions.Default(ctx)
		userRole := session.Get("userRole")
		// Return the updated program, passing data in a map
		ctx.HTML(http.StatusOK, "program-row.html", gin.H{
			"Item":     program,
			"UserRole": userRole,
		})
	}
}
