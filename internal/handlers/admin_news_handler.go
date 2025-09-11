package handler

import (
	"log"
	"net/http"

	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/DmytroPI-dev/clinic-golang/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Rendering news page“
func ShowNewsPage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var News []models.News
		db.Order("id asc").Find(&News)
		session := sessions.Default(c)
		username := session.Get("username")
		c.HTML(http.StatusOK, "news.html", gin.H{
			"Title": "Manage News",
			"User":  username,
			"Items": News,
		})
	}
}

func AdminShowNewsForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "news-form.html", gin.H{
		"News": models.News{},
	})
}

// Create new news template
func AdminCreateNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse form data from the request
		Title := ctx.PostForm("title")
		Header := ctx.PostForm("header")
		Description := ctx.PostForm("description")
		Features := ctx.PostForm("features")

		// Create a new news model instance with the data
		newNews := models.News{
			Title: Title, Header: Header, Description: Description, Features: Features,
			// Set translated fields to default language
			TitlePL:       Title,
			HeaderPL:      Header,
			DescriptionPL: Description,
			FeaturesPL:    Features,
			TitleEN:       Title,
			HeaderEN:      Header,
			DescriptionEN: Description,
			FeaturesEN:    Features,
			TitleUK:       Title,
			HeaderUK:      Header,
			DescriptionUK: Description,
			FeaturesUK:    Features,
		}

		// Process and save imageLeft if provided
		fileLeft, errLeft := ctx.FormFile("image_left")
		if errLeft == nil { // No error means file was provided
			savedPathLeft, err := utils.ProcessAndSaveImages(fileLeft)
			if err != nil {
				log.Printf("Failed to process and save imageLeft: %s", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			newNews.ImageLeft = savedPathLeft
		}

		// Process and save imageRight if provided
		fileRight, errRight := ctx.FormFile("image_right")
		if errRight == nil { // No error means file was provided
			savedPathRight, err := utils.ProcessAndSaveImages(fileRight)
			if err != nil {
				log.Printf("Failed to process and save imageRight: %s", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			newNews.ImageRight = savedPathRight
		}

		// Save the newly created news item to DB
		if err := db.Create(&newNews).Error; err != nil {
			log.Printf("Failed to create news: %s", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Render and return HTML fragment for new row
		ctx.HTML(http.StatusOK, "news-row.html", newNews)
	}
}

// AdminDeleteNews handles the deletion of News from the admin panel.
func AdminDeleteNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Delete the News from the database
		if err := db.Delete(&models.News{}, id).Error; err != nil {
			// We will just log the error now, later adding flash error
			log.Printf("Failed to delete News with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return an empty response
		ctx.String(http.StatusOK, "")
	}
}

// ShowEditNewsForm finds a News by ID and renders the edit form.
func AdminShowEditNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Find the news in the database
		var news models.News
		if err := db.First(&news, id).Error; err != nil {
			// Handle the case where no record found
			log.Printf("Failed to find news with ID %s: %s", id, err)
			ctx.Status(http.StatusNotFound)
			return
		}
		// Render the edit form with the news data
		ctx.HTML(http.StatusOK, "news-form.html", gin.H{
			"News": news,
		})
	}
}

// UpdateNews handles the submission of the edit news form.
func AdminUpdateNews(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")

		// Find existing news
		var news models.News
		if err := db.First(&news, id).Error; err != nil {
			// Handle the case where no record found
			log.Printf("Failed to find news with ID %s: %s", id, err)
			ctx.Status(http.StatusNotFound)
			return
		}

		//Parse data from the request
		news.Title = ctx.PostForm("title")
		news.Header = ctx.PostForm("header")
		news.Description = ctx.PostForm("description")
		news.Features = ctx.PostForm("features")
		// Will update translation fields later

		// Process and save imageLeft if provided
		fileLeft, errLeft := ctx.FormFile("image_left")
		if errLeft == nil { // No error means file was provided
			savedPathLeft, err := utils.ProcessAndSaveImages(fileLeft)
			if err != nil {
				log.Printf("Failed to process and save imageLeft: %s", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			news.ImageLeft = savedPathLeft
		}

		// Process and save imageRight if provided
		fileRight, errRight := ctx.FormFile("image_right")
		if errRight == nil { // No error means file was provided
			savedPathRight, err := utils.ProcessAndSaveImages(fileRight)
			if err != nil {
				log.Printf("Failed to process and save imageRight: %s", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			news.ImageRight = savedPathRight
		}

		// Save updates to the DB
		if err := db.Save(&news).Error; err != nil {
			log.Printf("Failed to update news with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return the updated news
		ctx.HTML(http.StatusOK, "news-row.html", news)
	}
}
