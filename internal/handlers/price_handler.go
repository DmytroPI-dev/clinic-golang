package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type PriceResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	ItemName   string    `json:"item_name"`
	Price      float32   `json:"price"`
	Category   string    `json:"category"`
	ItemNamePL string    `json:"item_name_pl"`
	ItemNameEN string    `json:"item_name_en"`
	ItemNameUK string    `json:"item_name_uk"`
}

func ListPrices(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var prices []models.Price
		// 1. Fetching all prices from the database.
		if err := db.Find(&prices).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Prices"})
			return
		}
		// 2. Mapping the database models to our API responce structs.
		var responces []PriceResponse
		for _, price := range prices {
			responces = append(responces, PriceResponse{
				ID:         price.ID,
				CreatedAt:  price.CreatedAt,
				ItemName:   price.ItemName,
				Price:      price.Price,
				Category:   price.Category,
				ItemNamePL: price.ItemNamePL,
				ItemNameEN: price.ItemNameEN,
				ItemNameUK: price.ItemNameUK,
			})
		}
		ctx.JSON(http.StatusOK, responces)
	}
}

func GetPrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from the URL parameter

		id := ctx.Param("id")
		var price models.Price

		// 2. Find first record matching ID.
		// Will use GORM `First` method for that
		if err := db.First(&price, id).Error; err != nil {
			// Handle the case where no record found.
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Price not found"})
				return
			} else {
				// Handle other database errors.
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Price"})
			}
			return
		}
		response := PriceResponse{
			ID:         price.ID,
			CreatedAt:  price.CreatedAt,
			ItemName:   price.ItemName,
			Price:      price.Price,
			Category:   price.Category,
			ItemNamePL: price.ItemNamePL,
			ItemNameEN: price.ItemNameEN,
			ItemNameUK: price.ItemNameUK,
		}
		// Sending the responce
		ctx.JSON(http.StatusOK, response)
	}
}

// CreatePriceRequest defines the structure for the request body when creating a price.
// We use `binding:"required"` for basic validation.
type CreatePriceRequest struct {
	ItemName string  `json:"item_name" binding:"required"`
	Price    float32 `json:"price" binding:"required"`
	Category string  `json:"category" binding:"required,len=2"`
}

func CreatePrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreatePriceRequest
		// 1. Bind the incoming JSON to the request struct.
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Create price instance
		price := models.Price{
			ItemName: request.ItemName,
			Price:    request.Price,
			Category: request.Category,
			// Set translated fields to default language
			ItemNamePL: request.ItemName,
			ItemNameEN: request.ItemName,
			ItemNameUK: request.ItemName,
		}
		// 2. Create price record in the database.
		if err := db.Create(&price).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Price"})
			return
		}
		// Return created record as a response
		// A 201 Created status will return
		ctx.JSON(http.StatusCreated, price)
	}
}

type UpdatePriceRequest struct {
	ItemName   string  `json:"item_name" binding:"required"`
	Price      float32 `json:"price" binding:"required"`
	Category   string  `json:"category" binding:"required,len=2"`
	ItemNamePL string  `json:"item_name_pl"`
	ItemNameEN string  `json:"item_name_en"`
	ItemNameUK string  `json:"item_name_uk"`
}

func UpdatePrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from URL
		id := ctx.Param("id")
		// 2. Find existing record in the database

		var price models.Price
		if err := db.First(&price, id).Error; err != nil {
			// Handle no record case
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Price not found"})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			}
			return
		}
		// 3. Bind the incoming JSON to a request struct.
		var request UpdatePriceRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 4.  Update the fields of the price model with the new data.
		price.ItemName = request.ItemName
		price.Price = request.Price
		price.Category = request.Category
		price.ItemNamePL = request.ItemNamePL
		price.ItemNameEN = request.ItemNameEN
		price.ItemNameUK = request.ItemNameUK
		// 5. Save the updated price in the database.
		if err := db.Save(&price).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Price"})
			return
		}
		// Return updated response
		ctx.JSON(http.StatusOK, price)
	}
}

func DeletePrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from URL
		id := ctx.Param("id")
		// 2. Find the price in the database
		result := db.Delete(&models.Price{}, id)
		// 3. Handle DB errors
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Price"})
			return
		}
		// 4. Check if record was deleted
		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Price not found"})
			return
		}
		// 5. Send success response.
		// The standard response for a successful DELETE is 204 No Content.
		ctx.Status(http.StatusNoContent)
	}
}
