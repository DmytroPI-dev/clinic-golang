package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"github.com/disintegration/imaging"
)

// ProcessAndSaveImage handles uploading, resizing, and saving an image.
// It returns the public path to the saved file or an error.
func ProcessAndSaveImages(file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Decode image
	img, err := imaging.Decode(src)
	if err != nil {
		return "", err
	}

	// Resizing image to max width 800px, preserving aspect ratio
	resized := imaging.Resize(img, 800, 0, imaging.Lanczos)

	// Create a unique name for the file to avoid collisions.
	uniqueFileName := fmt.Sprintf("%d%s", time.Now().Unix(), filepath.Base(file.Filename))

	// Define the filesystem path for saving the image.
	// This will be relative to the application's execution directory, e.g., "uploads/image.jpg"
	savePath := filepath.Join("./backend/uploads", uniqueFileName)

	// Save the resized image to the filesystem.
	err = imaging.Save(resized, savePath)
	if err != nil {
		return "", err
	}

	// Return the public URL path that the frontend will use, e.g., "/uploads/image.jpg"
	// We use filepath.ToSlash to ensure forward slashes are used in the URL, which is important for cross-platform compatibility.
	return filepath.ToSlash(filepath.Join("/", "uploads", uniqueFileName)), nil
}
