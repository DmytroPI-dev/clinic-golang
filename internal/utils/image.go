package utils

import (
	"fmt"
	"github.com/disintegration/imaging"
	"mime/multipart"
	"path/filepath"
	"time"
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

	// Create unique name
	uniqueFileName := fmt.Sprintf("%d%s", time.Now().Unix(), filepath.Base(file.Filename))
	savePath := filepath.Join("uploads", uniqueFileName)

	// Save resized image
	err = imaging.Save(resized, savePath)
	if err != nil {
		return "", err
	}

	return "/" + savePath, nil

}
