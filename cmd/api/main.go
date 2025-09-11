package main

import (
	"github.com/DmytroPI-dev/clinic-golang/internal/config"
	"github.com/DmytroPI-dev/clinic-golang/internal/database"
	"github.com/DmytroPI-dev/clinic-golang/internal/handler"
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
)

var funcMap = template.FuncMap{
	"Title": utils.Title,
	"Dict":  utils.Dict,
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	adminTpl := func(name string) string {
		return "templates/admin/" + name
	}

	layout := adminTpl("layout.html")
	programForm := adminTpl("program-form.html")
	programRow := adminTpl("program-row.html")
	priceForm := adminTpl("price-form.html")
	priceRow := adminTpl("price-row.html")

	r.AddFromFilesFuncs("programs.html", funcMap, layout,  adminTpl("programs.html"), programForm, programRow)
	r.AddFromFilesFuncs("prices.html", funcMap, layout,  adminTpl("prices.html"), priceForm, priceRow)

	// For HTMX partials and standalone pages
	partials := []string{
		"login.html",
		"program-form.html",
		"program-row.html",
		"price-form.html",
		"price-row.html",
	}
	for _, partial := range partials {
		r.AddFromFilesFuncs(partial, funcMap, adminTpl(partial))
	}

	return r
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
	ShowPage     func(*gorm.DB) gin.HandlerFunc
	ShowNewForm  gin.HandlerFunc
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

			// Admin CRUD pages
			registerAdminCrudRoutes(authenticated.Group("/programs"), db, AdminCrudHandlers{
				ShowPage:     handler.ShowProgramsPage,
				ShowNewForm:  handler.AdminShowNewProgramForm,
				Create:       handler.AdminCreateNewProgram,
				ShowEditForm: handler.AdminShowEditProgramForm,
				Update:       handler.AdminUpdateProgram,
				Delete:       handler.AdminDeleteProgram,
			})
			registerAdminCrudRoutes(authenticated.Group("/prices"), db, AdminCrudHandlers{
				ShowPage:     handler.ShowPricesPage,
				ShowNewForm:  handler.AdminShowNewPriceForm,
				Create:       handler.AdminCreateNewPrice,
				ShowEditForm: handler.AdminShowEditPriceForm,
				Update:       handler.AdminUpdatePrice,
				Delete:       handler.AdminDeletePrice,
			})
		}
	}

	//Testing route
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Start server
	serverAddress := "localhost:" + cfg.ServerPort
	log.Printf("Starting server on %s", serverAddress)
	router.Run(serverAddress)

}
