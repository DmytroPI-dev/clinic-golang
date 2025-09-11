package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)


// Rendering login page
func ShowLoginPage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	flashes := session.Flashes("error")
	// Important: Save the session to ensure flashes are cleared for the next request
	session.Save()
	renderData := gin.H{}
	if len(flashes) > 0 {
		renderData["error"] = flashes[0]
	}
	ctx.HTML(http.StatusOK, "login.html", renderData)
}

// Handle login
func HandleLogin(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session
		session := sessions.Default(ctx)
		// Obtain username and login
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		var user models.User
		if err := db.Where("user_name = ?", username).First(&user).Error; err != nil {
			// User not found message
			session.AddFlash("Invalid username or password", "error")
			session.Save()
			// Redirect to login
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}
		// Check password and hash
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			// Password do not match
			session.AddFlash("Invalid username or password", "error")
			session.Save()
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}
		// Create session
		session.Set("userID", user.ID)
		session.Set("username", user.UserName)
		session.Save()

		// Redirect to admin dashboard
		ctx.Redirect(http.StatusFound, "/admin/programs")
	}
}

// Auth Required middleware to check user auth
func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get("userID")
		// If user is not in the session, not logging user.
		if userID == nil {
			// Aborting request and redirecting to login page
			ctx.Abort()
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}
		ctx.Next()
	}
}

// HandleLogout clears the user's session and redirects to the login page.
func HandleLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusFound, "/admin/login")
}

