package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// ProgramResponse defines the structure of the JSON response for a program.
// We use `json:"..."` tags to control the field names in the JSON output.
type ProgramResponse struct {
	ID            uint   `json:"pk"`
	Title         string `json:"title"`
	TitleUK       string `json:"title_uk"`
	TitlePL       string `json:"title_pl"`
	TitleEN       string `json:"title_en"`
	Description   string `json:"description"`
	DescriptionUK string `json:"description_uk"`
	DescriptionPL string `json:"description_pl"`
	DescriptionEN string `json:"description_en"`
	Results       string `json:"results"`
	ResultsUK     string `json:"results_uk"`
	ResultsPL     string `json:"results_pl"`
	ResultsEN     string `json:"results_en"`
	Category      string `json:"category"`
}

// ListPrograms is the handler for fetching all programs.
// It accepts the GORM database connection as an argument.
func ListPrograms(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var programs []models.Program
		// 1. Fetching all programs from the database.
		if err := db.Find(&programs).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Programs"})
			return
		}

		// 2. Mapping the database models to our API responce structs.
		var responses []ProgramResponse
		for _, program := range programs {
			responses = append(responses, ProgramResponse{
				ID:            program.ID,
				Title:         program.Title,
				TitleUK:       program.TitleUK,
				TitlePL:       program.TitlePL,
				TitleEN:       program.TitleEN,
				Description:   program.Description,
				DescriptionUK: program.DescriptionUK,
				DescriptionPL: program.DescriptionPL,
				DescriptionEN: program.DescriptionEN,
				Results:       program.Results,
				ResultsUK:     program.ResultsUK,
				ResultsPL:     program.ResultsPL,
				ResultsEN:     program.ResultsEN,
				Category:      program.Category,
			})
		}
		ctx.JSON(http.StatusOK, responses)
	}
}

// GetProgram is the handler for fetching a single program by its ID.
func GetProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from the URL parameter

		id := ctx.Param("id")
		var program models.Program

		// 2. Find the first record that matches the ID.
		// Will use GORM `First` method for that
		if err := db.First(&program, id).Error; err != nil {
			// Handle the case where no record found.
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
				return
			} else {
				// Handle other database errors.
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch program"})
			}
			return
		}
		response := ProgramResponse{
			ID:            program.ID,
			Title:         program.Title,
			TitleUK:       program.TitleUK,
			TitlePL:       program.TitlePL,
			TitleEN:       program.TitleEN,
			Description:   program.Description,
			DescriptionUK: program.DescriptionUK,
			DescriptionPL: program.DescriptionPL,
			DescriptionEN: program.DescriptionEN,
			Results:       program.Results,
			ResultsUK:     program.ResultsUK,
			ResultsPL:     program.ResultsPL,
			ResultsEN:     program.ResultsEN,
			Category:      program.Category,
		}
		// Sending the responce
		ctx.JSON(http.StatusOK, response)
	}
}

// CreateProgramRequest defines the structure for the request body when creating a program.
// We use `binding:"required"` for basic validation.
type CreateProgramRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Results     string `json:"results"`
	Category    string `json:"category" binding:"required"`
}

// CreateProgram is the handler for creating a new program.
func CreateProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreateProgramRequest

		// 1. Bind the incoming JSON to the request struct.
		// If there's a validation error, it will be caught here.
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		program := models.Program{
			Title:       request.Title,
			Description: request.Description,
			Results:     request.Results,
			Category:    request.Category,
			// Set translated fields to default language
			TitlePL:       request.Title,
			TitleEN:       request.Title,
			TitleUK:       request.Title,
			DescriptionPL: request.Description,
			DescriptionEN: request.Description,
			DescriptionUK: request.Description,
			ResultsPL:     request.Results,
			ResultsEN:     request.Results,
			ResultsUK:     request.Results,
		}
		// 2. Create the program in the database.
		if err := db.Create(&program).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create program"})
			return
		}
		// Return created record as a response
		// A 201 Created status will return
		ctx.JSON(http.StatusCreated, program)
	}

}

// UpdateProgramRequest defines the structure for the request body when updating a program.

type UpdateProgramRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Results     string `json:"results"`
	Category    string `json:"category" binding:"required,len=2"`
	// Will add translated fields to allow them to be updated
	TitlePL       string `json:"title_pl"`
	TitleEN       string `json:"title_en"`
	TitleUK       string `json:"title_uk"`
	DescriptionPL string `json:"description_pl"`
	DescriptionEN string `json:"description_en"`
	DescriptionUK string `json:"description_uk"`
	ResultsPL     string `json:"results_pl"`
	ResultsEN     string `json:"results_en"`
	ResultsUK     string `json:"results_uk"`
}

func UpdateProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from URL
		id := ctx.Param("id")
		// 2. Find existing record in the database

		var program models.Program
		if err := db.First(&program, id).Error; err != nil {
			// Handle the case where no record found
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			}
			return
		}
		// 3. Bind the incoming JSON to a request struct.
		var request UpdateProgramRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 4.  Update the fields of the program model with the new data.
		program.Title = request.Title
		program.Description = request.Description
		program.Results = request.Results
		program.Category = request.Category
		program.TitlePL = request.TitlePL
		program.TitleEN = request.TitleEN
		program.TitleUK = request.TitleUK
		program.DescriptionPL = request.DescriptionPL
		program.DescriptionEN = request.DescriptionEN
		program.DescriptionUK = request.DescriptionUK
		program.ResultsPL = request.ResultsPL
		program.ResultsEN = request.ResultsEN
		program.ResultsUK = request.ResultsUK

		//5. Save the updated record to the database.
		if err := db.Save(&program).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update program"})
			return
		}

		// Return updated response
		ctx.JSON(http.StatusOK, program)
	}
}

// DeleteProgram is the handler for deleting a program.

func DeleteProgram(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from URL
		id := ctx.Param("id")
		// 2. Use GORM to delete the program using ID.
		// We execute the delete and check the number of rows affected.
		result := db.Delete(&models.Program{}, id)

		// 3. Handle DB errors
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete program"})
			return
		}
		// 4. Check if record was deleted
		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
			return
		}
		// 5. Send success response.
		// The standard response for a successful DELETE is 204 No Content.
		ctx.Status(http.StatusNoContent)
	}
}
