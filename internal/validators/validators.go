package validators

import (
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	uppercaseRegex   = regexp.MustCompile(`\p{Lu}`)
	lowercaseRegex   = regexp.MustCompile(`\p{Ll}`)
	numberRegex      = regexp.MustCompile(`\d`)
	specialCharRegex = regexp.MustCompile(`[@$!%*#?&]`)
)

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if !uppercaseRegex.MatchString(password) {
		return false
	}
	if !lowercaseRegex.MatchString(password) {
		return false
	}
	if !numberRegex.MatchString(password) {
		return false
	}
	if !specialCharRegex.MatchString(password) {
		return false
	}
	return true
}

func ValidateFileTypeExtension(file *multipart.FileHeader) bool {
	// Проверка на допустимый тип `Content-Type`
	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/svg+xml": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		return false
	}

	extension := strings.ToLower(filepath.Ext(file.Filename))
	// Проверка на допустимое расширение файла
	allowedExtensions := map[string]bool{
		".jpeg": true,
		".jpg":  true,
		".png":  true,
		".svg":  true,
	}
	if !allowedExtensions[extension] {
		return false
	}

	return true
}
