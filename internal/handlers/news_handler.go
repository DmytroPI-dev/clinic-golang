package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// NewsResponse defines the structure of the JSON response for a program.
// We use `json:"..."` tags to control the field names in the JSON output.
type NewsResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Header        string    `json:"header"`
	Description   string    `json:"description"`
	Features      string    `json:"features"`
	PostedOn      time.Time `json:"posted_on"`
	ImageLeft     string    `json:"image_left"`
	ImageRight    string    `json:"image_right"`
	TitlePL       string    `json:"title_pl"`
	DescriptionPL string    `json:"description_pl"`
	HeaderPL      string    `json:"header_pl"`
	FeaturesPL    string    `json:"features_pl"`
	TitleEN       string    `json:"title_en"`
	DescriptionEN string    `json:"description_en"`
	HeaderEN      string    `json:"header_en"`
	FeaturesEN    string    `json:"features_en"`
	TitleUK       string    `json:"title_uk"`
	DescriptionUK string    `json:"description_uk"`
	HeaderUK      string    `json:"header_uk"`
	FeaturesUK    string    `json:"features_uk"`
}

// ListNews is the handler for fetching all News.
// It accepts the GORM database connection as an argument.
func ListNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var news []models.News
		// 1. Fetching all News from the database.
		if err := db.Find(&news).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch News"})
			return
		}
		// 2. Mapping the database models to our API responce structs.
		var responces []NewsResponse
		for _, singleNews := range news {
			responces = append(responces, NewsResponse{
				ID:            singleNews.ID,
				CreatedAt:     singleNews.CreatedAt,
				Title:         singleNews.Title,
				Header:        singleNews.Header,
				Description:   singleNews.Description,
				Features:      singleNews.Features,
				PostedOn:      singleNews.PostedOn,
				ImageLeft:     singleNews.ImageLeft,
				ImageRight:    singleNews.ImageRight,
				TitlePL:       singleNews.TitlePL,
				DescriptionPL: singleNews.DescriptionPL,
				HeaderPL:      singleNews.HeaderPL,
				FeaturesPL:    singleNews.FeaturesPL,
				TitleEN:       singleNews.TitleEN,
				DescriptionEN: singleNews.DescriptionEN,
				HeaderEN:      singleNews.HeaderEN,
				FeaturesEN:    singleNews.FeaturesEN,
				TitleUK:       singleNews.TitleUK,
				DescriptionUK: singleNews.DescriptionUK,
				HeaderUK:      singleNews.HeaderUK,
				FeaturesUK:    singleNews.FeaturesUK,
			})
		}
		ctx.JSON(http.StatusOK, responces)
	}
}

func GetNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from the URL parameter

		id := ctx.Param("id")
		var news models.News

		// 2. Find the first record that matches the ID.
		// Will use GORM `First` method for that
		if err := db.First(&news, id).Error; err != nil {
			// Handle the case where no record found.
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
				return
			} else {
				// Handle other database errors.
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch News"})
			}
			return
		}
		response := NewsResponse{
			ID:            news.ID,
			CreatedAt:     news.CreatedAt,
			Title:         news.Title,
			Header:        news.Header,
			Description:   news.Description,
			Features:      news.Features,
			PostedOn:      news.PostedOn,
			ImageLeft:     news.ImageLeft,
			ImageRight:    news.ImageRight,
			TitlePL:       news.TitlePL,
			DescriptionPL: news.DescriptionPL,
			HeaderPL:      news.HeaderPL,
			FeaturesPL:    news.FeaturesPL,
			TitleEN:       news.TitleEN,
			DescriptionEN: news.DescriptionEN,
			HeaderEN:      news.HeaderEN,
			FeaturesEN:    news.FeaturesEN,
			TitleUK:       news.TitleUK,
			DescriptionUK: news.DescriptionUK,
			HeaderUK:      news.HeaderUK,
			FeaturesUK:    news.FeaturesUK,
		}
		ctx.JSON(http.StatusOK, response)
	}
}

type CreateNewsRequest struct {
	Title       string    `json:"title" binding:"required"`
	Header      string    `json:"header" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Features    string    `json:"features" binding:"required"`
	PostedOn    time.Time `json:"posted_on" binding:"required"`
	ImageLeft   string    `json:"image_left" binding:"required"`
	ImageRight  string    `json:"image_right" binding:"required"`
}

func CreateNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreateNewsRequest
		// 1.  Bind the incoming JSON to the request struct.
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Create News instance
		singleNews := models.News{
			Title:       request.Title,
			Header:      request.Header,
			Description: request.Description,
			Features:    request.Features,
			PostedOn:    request.PostedOn,
			ImageLeft:   request.ImageLeft,
			ImageRight:  request.ImageRight,
			// Set translated fields to default language
			TitlePL:       request.Title,
			HeaderPL:      request.Header,
			DescriptionPL: request.Description,
			FeaturesPL:    request.Features,
			TitleEN:       request.Title,
			HeaderEN:      request.Header,
			DescriptionEN: request.Description,
			FeaturesEN:    request.Features,
			TitleUK:       request.Title,
			HeaderUK:      request.Header,
			DescriptionUK: request.Description,
			FeaturesUK:    request.Features,
		}
		// 2. Create news record in the database.
		if err := db.Create(&singleNews).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create News"})
			return
		}
		// Return created record as a response
		// A 201 Created status will return
		ctx.JSON(http.StatusCreated, singleNews)
	}
}

type UpdateNewsRequest struct {
	Title         string    `json:"title" binding:"required"`
	Header        string    `json:"header" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	Features      string    `json:"features" binding:"required"`
	PostedOn      time.Time `json:"posted_on" binding:"required"`
	ImageLeft     string    `json:"image_left" binding:"required"`
	ImageRight    string    `json:"image_right" binding:"required"`
	TitlePL       string    `json:"title_pl"`
	HeaderPL      string    `json:"header_pl"`
	DescriptionPL string    `json:"description_pl"`
	FeaturesPL    string    `json:"features_pl"`
	TitleEN       string    `json:"title_en"`
	HeaderEN      string    `json:"header_en"`
	DescriptionEN string    `json:"description_en"`
	FeaturesEN    string    `json:"features_en"`
	TitleUK       string    `json:"title_uk"`
	HeaderUK      string    `json:"header_uk"`
	DescriptionUK string    `json:"description_uk"`
	FeaturesUK    string    `json:"features_uk"`
}

func UpdateNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from URL
		id := ctx.Param("id")
		// 2. Find existing	record
		var singleNews models.News
		if err := db.First(&singleNews, id).Error; err != nil {
			// Handle no record case
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			}
			return
		}
		// Binding incoming JSON to a request struct.
		var request UpdateNewsRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Update the fields of the News model with the new data
		singleNews.Title = request.Title
		singleNews.Header = request.Header
		singleNews.Description = request.Description
		singleNews.Features = request.Features
		singleNews.PostedOn = request.PostedOn
		singleNews.ImageLeft = request.ImageLeft
		singleNews.ImageRight = request.ImageRight
		singleNews.TitlePL = request.TitlePL
		singleNews.HeaderPL = request.HeaderPL
		singleNews.DescriptionPL = request.DescriptionPL
		singleNews.FeaturesPL = request.FeaturesPL
		singleNews.TitleEN = request.TitleEN
		singleNews.HeaderEN = request.HeaderEN
		singleNews.DescriptionEN = request.DescriptionEN
		singleNews.FeaturesEN = request.FeaturesEN
		singleNews.TitleUK = request.TitleUK
		singleNews.HeaderUK = request.HeaderUK
		singleNews.DescriptionUK = request.DescriptionUK
		singleNews.FeaturesUK = request.FeaturesUK

		// Saving updated news to database
		if err := db.Save(&singleNews).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update News"})
			return
		}
		// Return updated response
		ctx.JSON(http.StatusOK, singleNews)
	}
}

func DeleteNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the ID from URL
		id := ctx.Param("id")
		// 2. Find the news in the database
		result := db.Delete(&models.News{}, id)
		// 3. Handle DB errors
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete News"})
			return
		}
		// 4. Check if record was deleted
		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}
