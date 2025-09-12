package main

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/config"
	"github.com/DmytroPI-dev/clinic-golang/internal/database"
	handler "github.com/DmytroPI-dev/clinic-golang/internal/handlers"
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"github.com/DmytroPI-dev/clinic-golang/internal/utils"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var funcMap = template.FuncMap{
	"Title": utils.Title,
	"Dict":  utils.Dict,
}

func loadTemplates() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()
	// creating adminTpl var- subpath to
	adminTpl := func(name string) string {
		return "templates/admin/" + name
	}

	layout := adminTpl("layout.html")

	// Program
	programForm := adminTpl("program-form.html")
	programRow := adminTpl("program-row.html")
	// Price
	priceForm := adminTpl("price-form.html")
	priceRow := adminTpl("price-row.html")
	// News
	newsForm := adminTpl("news-form.html")
	newsRow := adminTpl("news-row.html")
	// Users
	usersRow := adminTpl("user-row.html")
	usersForm := adminTpl("user-form.html")
	// Access forbidden
	forbidden := adminTpl("403.html")

	// Configure HTML template rendering
	renderer.AddFromFilesFuncs("programs.html", funcMap, layout, adminTpl("programs.html"), programForm, programRow)
	renderer.AddFromFilesFuncs("prices.html", funcMap, layout, adminTpl("prices.html"), priceForm, priceRow)
	renderer.AddFromFilesFuncs("news.html", funcMap, layout, adminTpl("news.html"), newsForm, newsRow)
	renderer.AddFromFilesFuncs("users.html", funcMap, layout, adminTpl("users.html"), usersForm, usersRow)
	renderer.AddFromFilesFuncs("403.html", funcMap, layout, forbidden)

	// For HTMX partials and standalone pages
	partials := []string{
		"login.html",
		"program-form.html",
		"program-row.html",
		"price-form.html",
		"price-row.html",
		"news-form.html",
		"news-row.html",
		"user-row.html",
		"user-form.html",
	}
	for _, partial := range partials {
		renderer.AddFromFilesFuncs(partial, funcMap, adminTpl(partial))
	}
	return renderer
}

// CrudHandlers defines a set of handlers for a standard RESTful resource.
type CrudHandlers struct {
	List   func(*gorm.DB) gin.HandlerFunc
	Get    func(*gorm.DB) gin.HandlerFunc
	Create func(*gorm.DB) gin.HandlerFunc
	Update func(*gorm.DB) gin.HandlerFunc
	Delete func(*gorm.DB) gin.HandlerFunc
}

// registerCrudRoutes registers the standard CRUD endpoints for a resource.
func registerCrudRoutes(group *gin.RouterGroup, db *gorm.DB, handlers CrudHandlers) {
	group.GET("/", handlers.List(db))
	group.GET("/:id", handlers.Get(db))
	group.POST("/", handlers.Create(db))
	group.PUT("/:id", handlers.Update(db))
	group.DELETE("/:id", handlers.Delete(db))
}

// AdminCrudHandlers defines a set of handlers for an admin panel resource.
type AdminCrudHandlers struct {
	ShowNewForm  gin.HandlerFunc
	ShowPage     func(*gorm.DB) gin.HandlerFunc
	Create       func(*gorm.DB) gin.HandlerFunc
	ShowEditForm func(*gorm.DB) gin.HandlerFunc
	Update       func(*gorm.DB) gin.HandlerFunc
	Delete       func(*gorm.DB) gin.HandlerFunc
}

// registerAdminCrudRoutes registers the admin CRUD endpoints for a resource.
func registerAdminCrudRoutes(group *gin.RouterGroup, db *gorm.DB, handlers AdminCrudHandlers) {
	group.GET("/", handlers.ShowPage(db))
	group.GET("/new", handlers.ShowNewForm)
	group.POST("/", handlers.Create(db))
	group.GET("/edit/:id", handlers.ShowEditForm(db))
	group.PUT("/:id", handlers.Update(db))
	group.DELETE("/:id", handlers.Delete(db))
}

func main() {
	//Load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load environment variables: %s", err)
	}

	//Connect to DB
	db, err := database.DB_Connect(cfg.DB_DSN)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	log.Println("Successfully connected to database")
	// Migrating data
	log.Println("Starting DB migration....")
	if err := db.AutoMigrate(&models.Program{}, &models.Price{}, &models.News{}, &models.User{}); err != nil {
		log.Fatalf("migration for models.Program failed: %s", err)
	}
	log.Println("Migration successful")

	// Creating Gin router
	router := gin.Default()
	router.Static("/static", "./static")
	// Serve frontend static files from the 'web/static' directory under a unique path
	router.Static("/ui-assets", "./web/static")

	router.SetFuncMap(funcMap)
	// Setting up session store
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	router.Use(sessions.Sessions("session", store))

	// Loading templates
	router.HTMLRender = loadTemplates()

	// Grouping API routes
	v1 := router.Group("/api/v1")
	{
		// API CRUD endpoints
		registerCrudRoutes(v1.Group("/programs"), db, CrudHandlers{
			List:   handler.ListPrograms,
			Get:    handler.GetProgram,
			Create: handler.CreateProgram,
			Update: handler.UpdateProgram,
			Delete: handler.DeleteProgram,
		})
		registerCrudRoutes(v1.Group("/prices"), db, CrudHandlers{
			List:   handler.ListPrices,
			Get:    handler.GetPrice,
			Create: handler.CreatePrice,
			Update: handler.UpdatePrice,
			Delete: handler.DeletePrice,
		})
		registerCrudRoutes(v1.Group("/news"), db, CrudHandlers{
			List:   handler.ListNews,
			Get:    handler.GetNews,
			Create: handler.CreateNews,
			Update: handler.UpdateNews,
			Delete: handler.DeleteNews,
		})
	}

	// Admin routes
	adminRoutes := router.Group("/admin")
	{
		// Public routes that don't require authentication
		adminRoutes.GET("/login", handler.ShowLoginPage)
		adminRoutes.POST("/login", handler.HandleLogin(db))

		// Authenticated routes
		authenticated := adminRoutes.Group("/")
		authenticated.Use(handler.AuthRequired())
		{
			authenticated.GET("/logout", handler.HandleLogout)
			authenticated.GET("/", func(c *gin.Context) {
				c.Redirect(http.StatusFound, "/admin/programs")
			})

			// Admin CRUD pages with Role-Based Access Control
			// Programs: Readers can view, Editors/Admins can modify.
			programsGroup := authenticated.Group("/programs")
			programsGroup.GET("/", handler.RoleRequired(models.Admin, models.Editor, models.Reader), handler.ShowProgramsPage(db))
			programsGroup.Use(handler.RoleRequired(models.Admin, models.Editor))
			{
				programsGroup.GET("/new", handler.AdminShowNewProgramForm)
				programsGroup.POST("/", handler.AdminCreateNewProgram(db))
				programsGroup.GET("/edit/:id", handler.AdminShowEditProgramForm(db))
				programsGroup.PUT("/:id", handler.AdminUpdateProgram(db))
				programsGroup.DELETE("/:id", handler.AdminDeleteProgram(db))
			}

			// Prices: Readers can view, Editors/Admins can modify.
			pricesGroup := authenticated.Group("/prices")
			pricesGroup.GET("/", handler.RoleRequired(models.Admin, models.Editor, models.Reader), handler.ShowPricesPage(db))
			pricesGroup.Use(handler.RoleRequired(models.Admin, models.Editor))
			{
				pricesGroup.GET("/new", handler.AdminShowNewPriceForm)
				pricesGroup.POST("/", handler.AdminCreateNewPrice(db))
				pricesGroup.GET("/edit/:id", handler.AdminShowEditPriceForm(db))
				pricesGroup.PUT("/:id", handler.AdminUpdatePrice(db))
				pricesGroup.DELETE("/:id", handler.AdminDeletePrice(db))
			}

			// News: Readers can view, Editors/Admins can modify.
			newsGroup := authenticated.Group("/news")
			newsGroup.GET("/", handler.RoleRequired(models.Admin, models.Editor, models.Reader), handler.ShowNewsPage(db))
			newsGroup.Use(handler.RoleRequired(models.Admin, models.Editor))
			{
				newsGroup.GET("/new", handler.AdminShowNewsForm)
				newsGroup.POST("/", handler.AdminCreateNews(db))
				newsGroup.GET("/edit/:id", handler.AdminShowEditNews(db))
				newsGroup.PUT("/:id", handler.AdminUpdateNews(db))
				newsGroup.DELETE("/:id", handler.AdminDeleteNews(db))
			}

			// Users: Only Admins can manage users.
			usersGroup := authenticated.Group("/users")
			usersGroup.Use(handler.RoleRequired(models.Admin))
			registerAdminCrudRoutes(usersGroup, db, AdminCrudHandlers{
				ShowPage:     handler.ShowUserPage,
				ShowNewForm:  handler.AdminShowNewUserForm,
				Create:       handler.AdminCreateUser,
				ShowEditForm: handler.AdminShowEditUserForm,
				Update:       handler.AdminUpdateUser,
				Delete:       handler.AdminDeleteUser,
			})
		}

		//Testing route
		router.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		router.NoRoute(func(c *gin.Context) {
			// For API routes that are not found, we want to return a JSON 404.
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
				return
			}

			// For any other route, we assume it's for the frontend.
			// This logic serves static files (like CSS, JS) if they have an extension,
			// and serves 'index.html' for paths without an extension, which is
			// a common pattern for Single-Page Applications (SPAs).
			dir, file := filepath.Split(c.Request.URL.Path)
			ext := filepath.Ext(file)

			if file == "" || ext == "" {
				// When a directory or a path without an extension is requested,
				// serve the main 'index.html' file from the 'web' directory.
				c.File("./web/templates/404.html")
			} else {
				// When a file with an extension is requested, serve it from the 'web' directory.
				c.File(filepath.Join("./web", dir, file))
			}
		})

		// Start server
		serverAddress := "localhost:" + cfg.ServerPort
		log.Printf("Starting server on %s", serverAddress)
		router.Run(serverAddress)

	}
}
