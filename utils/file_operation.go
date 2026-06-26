package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize       = 10 * 1024 * 1024 // 10 MB
	UploadRoots       = "uploads"
	TaskAttachmentDir = "/user/task_attachment/"
	ProfilePictureDir = "/user/profile_picture/"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".pdf":  true,
}
var alloweImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}
var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
}
var allowedImageMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
}

func SaveTaskAttachment(file *multipart.FileHeader, context *gin.Context) (*string, error) {
	// File Size Validation
	if file.Size > MaxFileSize {
		return nil, errors.New("File size cannot exceed 10 MB!")
	}

	// File Extension Type Validation
	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[extension] {
		return nil, errors.New("Invalid file extension!")
	}

	// File MIME Type Validation
	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	n, err := openedFile.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	contentType := http.DetectContentType(buffer[:n])
	if !allowedMimeTypes[contentType] {
		// Special case for PDF
		if extension == ".pdf" && contentType == "application/octet-stream" {
			// Accept it
		} else {
			return nil, errors.New("Invalid content type: " + contentType)
		}
	}

	// Save the Attachment
	filename, err := SaveFile(file, TaskAttachmentDir, extension, context)
	if err != nil {
		return nil, err
	}

	return filename, nil
}

func SaveProfilePicture(file *multipart.FileHeader, context *gin.Context) (*string, error) {
	// Validate file size
	if file.Size > MaxFileSize {
		return nil, errors.New("Image size exceed 10 MB")
	}

	// Validate file extension
	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !alloweImageExtensions[extension] {
		return nil, errors.New("File is not supported")
	}

	// Validate MIME type
	fileOpened, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileOpened.Close()

	buffer := make([]byte, 512)
	n, err := fileOpened.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	contentType := http.DetectContentType(buffer[:n])
	if !allowedImageMimeTypes[contentType] {
		return nil, errors.New("Invalid content type: " + contentType)
	}

	// Save profile picture
	filename, err := SaveFile(file, ProfilePictureDir, extension, context)
	if err != nil {
		return nil, err
	}

	return filename, nil
}

func SaveFile(file *multipart.FileHeader, directory, extension string, context *gin.Context) (*string, error) {
	// Create directory
	u_id := strconv.FormatInt(context.GetInt64("user_id"), 10)
	os.MkdirAll(UploadRoots+directory+u_id, os.ModePerm)

	// Upload the file
	filename := fmt.Sprintf("task_%d_user_%s%s", time.Now().UnixNano(), u_id, extension)
	path := filepath.Join(UploadRoots, directory, u_id, filename)

	err := context.SaveUploadedFile(file, path)
	if err != nil {
		return nil, err
	}

	return &filename, nil
}

func RemoveFileAttachment(fileName *string, dir string, u_id int64) error {
	path := UploadRoots + dir + strconv.FormatInt(u_id, 10) + "/" + *fileName

	_, err := os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
