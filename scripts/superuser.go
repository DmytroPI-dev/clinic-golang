package main

// usage go run ./scripts/superuser -user=new_admin_username -password=a_very_strong_password -email=your@email.com

import (
	"errors"
	"flag"
	"github.com/DmytroPI-dev/clinic-golang/internal/config"
	"github.com/DmytroPI-dev/clinic-golang/internal/database"
	"github.com/DmytroPI-dev/clinic-golang/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Define and parse CLI flags for superuser, password, email
	username := flag.String("user", "", "Username for the superuser")
	password := flag.String("password", "", "Password for the superuser")
	email := flag.String("email", "", "Email for the superuser")
	flag.Parse()

	if *username == "" || *password == "" || *email == "" {
		log.Fatal("Superusername, password, and email are required")
	}

	// Loading config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load environment variables: %s", err)
	}

	// Connect to DB
	db, err := database.DB_Connect(cfg.DB_DSN)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	log.Println("Successfully connected to database")

	// Check if user exists
	var existingUser models.User
	err = db.Where("user_name = ?", *username).First(&existingUser).Error
	if err == nil {
		log.Fatalf("Superuser with username '%s' already exists, try another username!", *username)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("Error checking for existing superuser: %s", err)
	}

	// Hashing the password and creating new superuser
	log.Println("Creating new superuser...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
	if err != nil {
		log.Fatalf("Could not hash password: %s", err)
	}

	// Create new admin user
	adminUser := models.User{
		UserName:     *username,
		PasswordHash: string(hashedPassword),
		Email:        *email,
		Role:         cfg.AdminRole,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatalf("Could not create superuser: %s", err)
	}
	log.Println("Superuser created successfully")
}
