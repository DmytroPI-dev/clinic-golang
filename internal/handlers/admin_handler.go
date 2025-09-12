package handler

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
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
		// Obtain userName and login
		userName := ctx.PostForm("userName")
		password := ctx.PostForm("password")

		var user models.User
		if err := db.Where("user_name = ?", userName).First(&user).Error; err != nil {
			// User not found message
			session.AddFlash("Invalid user name or password", "error")
			session.Save()
			// Redirect to login
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}
		// Check password and hash
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			// Password do not match
			session.AddFlash("Invalid user name or password", "error")
			session.Save()
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}
		// Create session
		session.Set("userID", user.ID)
		session.Set("userName", user.UserName)
		session.Set("userRole", user.Role)
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

// RoleRequired is a middleware to check if the user has one of the allowed roles.
func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userRoleVal := session.Get("userRole")
		// Check if userRole exists in the session
		if userRoleVal == nil {
			// This should ideally be caught by AuthRequired middleware first.
			ctx.Redirect(http.StatusFound, "/admin/login") // Not logged in
			ctx.Abort()
			return
		}

		userRole, ok := userRoleVal.(string)
		if !ok {
			// The userRole in session is not a string, which is unexpected.
			log.Printf("userRole in session is not a string: %v", userRoleVal)
			ctx.HTML(http.StatusForbidden, "403.html", gin.H{"Title": "Forbidden"})
			ctx.Abort()
			return
		}

		// Check if the user's role is in the list of allowed roles
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			// User's role is not permitted. Show a "Forbidden" error.
			ctx.HTML(http.StatusForbidden, "403.html", gin.H{"Title": "Forbidden"})
			ctx.Abort()
			return
		}
		// Role is permitted, continue to the handler
		ctx.Next()
	}
}
