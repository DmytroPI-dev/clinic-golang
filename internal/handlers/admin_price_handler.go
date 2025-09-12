package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Rendering price page
func ShowPricesPage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var prices []models.Price
		db.Order("id asc").Find(&prices)
		session := sessions.Default(c)
		userName := session.Get("userName")
		userRole := session.Get("userRole")
		c.HTML(http.StatusOK, "prices.html", gin.H{
			"Title":    "Manage Prices",
			"User":     userName,
			"UserRole": userRole,
			"Items":    prices,
		})
	}
}

func AdminShowNewPriceForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "price-form.html", gin.H{
		"Categories": models.AllCategories,
		"Price":      models.Price{},
	})
}

// Create new price template
func AdminCreateNewPrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse form data from the request
		itemName := ctx.PostForm("itemName")
		category := ctx.PostForm("category")
		priceStr := ctx.PostForm("price")
		price, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			log.Printf("Failed to parse price: %s", err)
			ctx.Status(http.StatusBadRequest)
			return
		}

		// Create a new price model instance with the data
		newPrice := models.Price{
			ItemName:   itemName,
			Price:      float32(price),
			Category:   category,
			ItemNamePL: itemName,
			ItemNameEN: itemName,
			ItemNameUK: itemName,
		}
		// Save the newly created price item to DB
		if err := db.Create(&newPrice).Error; err != nil {
			log.Printf("Failed to create price: %s", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Render and return HTML fragment for new row
		ctx.HTML(http.StatusOK, "price-row.html", newPrice)
	}
}

// AdminDeletePrice handles the deletion of a price from the admin panel.
func AdminDeletePrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Delete the price from the database
		if err := db.Delete(&models.Price{}, id).Error; err != nil {
			// We will just log the error now, later adding flash error
			log.Printf("Failed to delete price with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return an empty response
		ctx.String(http.StatusOK, "")
	}
}

// ShowEditPriceForm finds a price by ID and renders the edit form.
func AdminShowEditPriceForm(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Find the price in the database
		var price models.Price
		if err := db.First(&price, id).Error; err != nil {
			// Handle the case where no record found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.HTML(http.StatusNotFound, "404.html", gin.H{"Title": "Not Found"})
			} else {
				log.Printf("Failed to find price with ID %s: %s", id, err)
				ctx.Status(http.StatusNotFound)
			}
			return
		}
		// Render the edit form with the price data
		ctx.HTML(http.StatusOK, "price-form.html", gin.H{
			"Categories": models.AllCategories,
			"Price":      price,
		})
	}
}

// UpdatePrice handles the submission of the edit price form.
func AdminUpdatePrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")

		// Find existing price
		var price models.Price
		if err := db.First(&price, id).Error; err != nil {
			// Handle the case where no record found
			log.Printf("Failed to find price with ID %s: %s", id, err)
			ctx.Status(http.StatusNotFound)
			return
		}

		//Parse data from the request
		price.ItemName = ctx.PostForm("itemName")
		priceStr := ctx.PostForm("price")
		if priceFloat, err := strconv.ParseFloat(priceStr, 32); err != nil {
			log.Printf("Failed to parse price '%s': %s", priceStr, err)
			ctx.Status(http.StatusBadRequest)
			return
		} else {
			price.Price = float32(priceFloat)
		}
		price.Category = ctx.PostForm("category")
		// Will update translation fields later

		// Save updates to the DB
		if err := db.Save(&price).Error; err != nil {
			log.Printf("Failed to update price with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return the updated price
		ctx.HTML(http.StatusOK, "price-row.html", price)
	}
}
