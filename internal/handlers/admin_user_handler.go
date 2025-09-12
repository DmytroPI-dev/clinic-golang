package handler

import (
	"errors"
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// Render users page
func ShowUserPage(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var users []models.User
		// Fetch all users data, excluding password hash field
		if err := db.Select("id", "user_name", "email", "role").Order("id asc").Find(&users).Error; err != nil {
			log.Printf("Failed to fetch users: %s", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Get session data
		session := sessions.Default(ctx)
		userName := session.Get("userName")
		userRole := session.Get("userRole")
		flashes := session.Flashes("error")
		if err := session.Save(); err != nil {
			log.Printf("Failed to save session to clear flashes: %s", err)
		}

		renderData := gin.H{
			"Title":    "Manage Users",
			"User":     userName,
			"UserRole": userRole,
			"Items":    users,
		}
		if len(flashes) > 0 {
			renderData["error"] = flashes[0]
		}
		// Render template
		ctx.HTML(http.StatusOK, "users.html", renderData)
	}
}

// Render new user page
func AdminShowNewUserForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user-form.html", gin.H{
		"Roles": models.AllRoles,
	})
}

// Create new user
func AdminCreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse form data from the request
		userName := ctx.PostForm("userName")
		email := ctx.PostForm("email")
		role := ctx.PostForm("role")
		password := ctx.PostForm("password")

		// Validate password
		if password == "" {
			ctx.HTML(http.StatusBadRequest, "user-form.html", gin.H{
				"Roles": models.AllRoles,
				"Error": "Password is required",
			})
			return
		}

		//Hash pasword
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			log.Printf("Password hash failed: %s", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// Create a new user model instance with data
		newUser := models.User{
			UserName:     userName,
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         role,
		}

		// Save new user
		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("Failed to create new user: %s", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Render and return HTML fragment for new row
		ctx.HTML(http.StatusOK, "user-row.html", newUser)
	}
}

// AdminDeleteUser handles the deletion of the user from the admin panel and database
func AdminDeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")

		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.Status(http.StatusNotFound)
				return
			}
			log.Printf("Failed to find User with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if user.Role == models.Admin {
			var adminCount int64
			if err := db.Model(&models.User{}).Where("role = ?", models.Admin).Count(&adminCount).Error; err != nil {
				log.Printf("Failed to count admin users: %s", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}

			if adminCount <= 1 {
				session := sessions.Default(ctx)
				session.AddFlash("Cannot delete the last admin user.", "error")
				if err := session.Save(); err != nil {
					log.Printf("Failed to save session: %s", err)
					ctx.Status(http.StatusInternalServerError)
					return
				}
				// Tell HTMX to refresh the page to show the flash message
				ctx.Header("HX-Refresh", "true")
				// Return a conflict status to indicate the nature of the error.
				ctx.Status(http.StatusConflict)
				return
			}
		}

		// Delete User from the database completely.
		if err := db.Unscoped().Delete(&models.User{}, id).Error; err != nil {
			log.Printf("Failed to delete User with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return an empty response
		ctx.String(http.StatusOK, "")
	}
}

// AdminShowEditUserForm finds a User by ID and renders the edit form.
func AdminShowEditUserForm(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Find the news in the database
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			// Handle the case where no record found
			log.Printf("Failed to find User with ID %s: %s", id, err)
			ctx.Status(http.StatusNotFound)
			return
		}
		// Render the edit form with the news data
		ctx.HTML(http.StatusOK, "user-form.html", gin.H{
			"User":  user,
			"Roles": models.AllRoles,
		})
	}
}

// UpdateUser handles the submission of the edit program form.
func AdminUpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from the URL
		id := ctx.Param("id")
		// Find existing user
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			// Handle the case where no record found
			log.Printf("Failed to find User with ID %s: %s", id, err)
			ctx.Status(http.StatusNotFound)
			return
		}
		// Parse form data from the request
		user.UserName = ctx.PostForm("userName")
		user.Email = ctx.PostForm("email")
		user.Role = ctx.PostForm("role")

		// Check if a new password was provided
		newPassword := ctx.PostForm("password")
		if newPassword != "" {
			// Hash the new password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
			if err != nil {
				log.Printf("Password hash failed: %s", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			user.PasswordHash = string(hashedPassword)
		}

		// Save updates to the DB
		if err := db.Save(&user).Error; err != nil {
			log.Printf("Failed to update User with ID %s: %s", id, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Return the updated user
		ctx.HTML(http.StatusOK, "user-row.html", user)
	}
}
