package handlers

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Images struct{}

func NewImages() *Images {
	return &Images{}
}

func (i *Images) GetImageAsBase64(filePath string) (string, error) {
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Determine mime type based on file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	mimeType := "image/jpeg"
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	}

	base64Data := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), nil
}

func (i *Images) GetVideoAsBase64(filePath string) (string, error) {
	videoData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Determine mime type based on file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	mimeType := "video/mp4"
	switch ext {
	case ".mp4":
		mimeType = "video/mp4"
	case ".webm":
		mimeType = "video/webm"
	case ".ogg":
		mimeType = "video/ogg"
	case ".mov":
		mimeType = "video/quicktime"
	case ".avi":
		mimeType = "video/x-msvideo"
	}

	base64Data := base64.StdEncoding.EncodeToString(videoData)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), nil
}
