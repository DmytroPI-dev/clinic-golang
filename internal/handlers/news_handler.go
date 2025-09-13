package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewsResponse defines the structure of the JSON response for a news.
// We use `json:"..."` tags to control the field names in the JSON output.
// NewsResponse matches a single item in the "results" array.
type NewsResponse struct {
	ID            uint      `json:"pk"`
	Title         string    `json:"title"`
	TitleUK       string    `json:"title_uk"`
	TitlePL       string    `json:"title_pl"`
	TitleEN       string    `json:"title_en"`
	Description   string    `json:"description"`
	DescriptionUK string    `json:"description_uk"`
	DescriptionPL string    `json:"description_pl"`
	DescriptionEN string    `json:"description_en"`
	Header        string    `json:"header"`
	HeaderUK      string    `json:"header_uk"`
	HeaderPL      string    `json:"header_pl"`
	HeaderEN      string    `json:"header_en"`
	Features      string    `json:"features"`
	FeaturesUK    string    `json:"features_uk"`
	FeaturesPL    string    `json:"features_pl"`
	FeaturesEN    string    `json:"features_en"`
	PostedOn      time.Time `json:"posted_on"`
	// Pointers are used for fields that can be null
	ImageLeft  *string `json:"image_left,omitempty"`
	ImageRight *string `json:"image_right,omitempty"`
}

// PaginatedNewsResponse matches the top-level paginated Django structure.
type PaginatedNewsResponse struct {
	Count    int64          `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []NewsResponse `json:"results"`
}

// ListNews is the handler for fetching all News.
// It accepts the GORM database connection as an argument.
func ListNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get pagination parameters from the query string (e.g., ?limit=10&page=1)
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "1"))
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		offset := (page - 1) * limit
		// Get total number of News
		var count int64
		if err := db.Model(&models.News{}).Count(&count).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count News"})
			return
		}
		// Fetching paginated list of News from the database.
		var newsItems []models.News
		if err := db.Limit(limit).Offset(offset).Order("posted_on desc").Find(&newsItems).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch News"})
			return
		}
		// Mapping the database models to our responce structs.
		results := make([]NewsResponse, 0, len(newsItems))
		for _, item := range newsItems { // This was a bug, appending to a pre-sized slice.
			// Handle nullable image fields
			var imgLeft, imgRight *string
			if item.ImageLeft != "" {
				imgLeft = &item.ImageLeft
			}
			if item.ImageRight != "" {
				imgRight = &item.ImageRight
			}

			results = append(results, NewsResponse{
				ID:            item.ID,
				Title:         item.Title,
				Header:        item.Header,
				Description:   item.Description,
				Features:      item.Features,
				TitlePL:       item.TitlePL,
				DescriptionPL: item.DescriptionPL,
				HeaderPL:      item.HeaderPL,
				FeaturesPL:    item.FeaturesPL,
				TitleEN:       item.TitleEN,
				DescriptionEN: item.DescriptionEN,
				HeaderEN:      item.HeaderEN,
				FeaturesEN:    item.FeaturesEN,
				TitleUK:       item.TitleUK,
				DescriptionUK: item.DescriptionUK,
				HeaderUK:      item.HeaderUK,
				FeaturesUK:    item.FeaturesUK,
				PostedOn:      item.PostedOn,
				ImageLeft:     imgLeft,
				ImageRight:    imgRight,
			})
		}
		// Build paginated response object
		var nextURL, prevURL *string
		baseURL := fmt.Sprintf("/api/news/?limit=%d", limit)

		if int64(page)*int64(limit) < count {
			url := fmt.Sprintf("%s&page=%d", baseURL, page+1)
			nextURL = &url
		}
		if page > 1 {
			url := fmt.Sprintf("%s&page=%d", baseURL, page-1)
			prevURL = &url
		}

		response := PaginatedNewsResponse{
			Count:    count,
			Next:     nextURL,
			Previous: prevURL,
			Results:  results,
		}

		ctx.JSON(http.StatusOK, response)
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
		// Handle nullable image fields
		var imgLeft, imgRight *string
		if news.ImageLeft != "" {
			imgLeft = &news.ImageLeft
		}
		if news.ImageRight != "" {
			imgRight = &news.ImageRight
		}

		response := NewsResponse{
			ID:            news.ID,
			Title:         news.Title,
			Header:        news.Header,
			Description:   news.Description,
			Features:      news.Features,
			PostedOn:      news.PostedOn,
			ImageLeft:     imgLeft,
			ImageRight:    imgRight,
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
		var imgLeft, imgRight *string
		if singleNews.ImageLeft != "" {
			imgLeft = &singleNews.ImageLeft
		}
		if singleNews.ImageRight != "" {
			imgRight = &singleNews.ImageRight
		}
		response := NewsResponse{
			ID:            singleNews.ID,
			Title:         singleNews.Title,
			Header:        singleNews.Header,
			Description:   singleNews.Description,
			Features:      singleNews.Features,
			PostedOn:      singleNews.PostedOn,
			ImageLeft:     imgLeft,
			ImageRight:    imgRight,
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
		}
		ctx.JSON(http.StatusCreated, response)
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
		var newsItem models.News
		if err := db.First(&newsItem, id).Error; err != nil {
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
		newsItem.Title = request.Title
		newsItem.Header = request.Header
		newsItem.Description = request.Description
		newsItem.Features = request.Features
		newsItem.PostedOn = request.PostedOn
		newsItem.ImageLeft = request.ImageLeft
		newsItem.ImageRight = request.ImageRight
		newsItem.TitlePL = request.TitlePL
		newsItem.HeaderPL = request.HeaderPL
		newsItem.DescriptionPL = request.DescriptionPL
		newsItem.FeaturesPL = request.FeaturesPL
		newsItem.TitleEN = request.TitleEN
		newsItem.HeaderEN = request.HeaderEN
		newsItem.DescriptionEN = request.DescriptionEN
		newsItem.FeaturesEN = request.FeaturesEN
		newsItem.TitleUK = request.TitleUK
		newsItem.HeaderUK = request.HeaderUK
		newsItem.DescriptionUK = request.DescriptionUK
		newsItem.FeaturesUK = request.FeaturesUK

		// Saving updated news to database
		if err := db.Save(&newsItem).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update News"})
			return
		}
		// Return updated response
		var imgLeft, imgRight *string
		if newsItem.ImageLeft != "" {
			imgLeft = &newsItem.ImageLeft
		}
		if newsItem.ImageRight != "" {
			imgRight = &newsItem.ImageRight
		}
		response := NewsResponse{
			ID:            newsItem.ID,
			Title:         newsItem.Title,
			Header:        newsItem.Header,
			Description:   newsItem.Description,
			Features:      newsItem.Features,
			PostedOn:      newsItem.PostedOn,
			ImageLeft:     imgLeft,
			ImageRight:    imgRight,
			TitlePL:       newsItem.TitlePL,
			DescriptionPL: newsItem.DescriptionPL,
			HeaderPL:      newsItem.HeaderPL,
			FeaturesPL:    newsItem.FeaturesPL,
			TitleEN:       newsItem.TitleEN,
			DescriptionEN: newsItem.DescriptionEN,
			HeaderEN:      newsItem.HeaderEN,
			FeaturesEN:    newsItem.FeaturesEN,
			TitleUK:       newsItem.TitleUK,
			DescriptionUK: newsItem.DescriptionUK,
			HeaderUK:      newsItem.HeaderUK,
			FeaturesUK:    newsItem.FeaturesUK,
		}
		ctx.JSON(http.StatusOK, response)
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
