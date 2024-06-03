package validators

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strings"
)

func getJSONTag(obj interface{}, fieldName string) string {
	t := reflect.TypeOf(obj).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == fieldName {
			return field.Tag.Get("json")
		}
		if field.Type.Kind() == reflect.Struct {
			jsonTag := getJSONTag(reflect.New(field.Type).Interface(), fieldName)
			if jsonTag != "" {
				return jsonTag
			}
		}
	}
	return ""
}

func CustomErrorMessage(err error, obj interface{}) string {
	var sb strings.Builder

	// Проверяем, является ли ошибка ошибкой валидации
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			field, _ := reflect.TypeOf(obj).Elem().FieldByName(e.StructField())
			jsonTag := getJSONTag(obj, e.StructField())
			if jsonTag == "" {
				jsonTag = e.Field()
			}

			var message string
			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("Поле `%s` является обязательным.", jsonTag)
			case "min":
				if field.Type.Kind() == reflect.String || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.String) {
					message = fmt.Sprintf("Поле `%s` должно быть не менее %s символов.", jsonTag, e.Param())
				} else {
					message = fmt.Sprintf("Поле `%s` должно быть не менее %s.", jsonTag, e.Param())
				}
			case "max":
				if field.Type.Kind() == reflect.String || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.String) {
					message = fmt.Sprintf("Поле `%s` должно быть не более %s символов.", jsonTag, e.Param())
				} else {
					message = fmt.Sprintf("Поле `%s` должно быть не более %s.", jsonTag, e.Param())
				}
			case "email":
				message = fmt.Sprintf("Поле `%s` должно быть корректным адресом электронной почты.", jsonTag)
			case "url":
				message = fmt.Sprintf("Поле `%s` должно быть корректным URL.", jsonTag)
			case "telegram":
				message = fmt.Sprintf("Поле `%s` должно быть корректным логином Telegram.", jsonTag)
			case "timeorder":
				message = fmt.Sprintf("Поле `%s` является некорректным. Дата начала события должна быть меньше даты конца события.", jsonTag)
			case "peopleorder":
				message = fmt.Sprintf("Поле `%s` является некорректным. Количество участников `от` должно быть меньше количества участников `до`.", jsonTag)
			default:
				message = fmt.Sprintf("Поле `%s` является некорректным.", jsonTag)
			}

			sb.WriteString(message)
			break // Возвращаем только первую ошибку
		}
	} else {
		sb.WriteString(err.Error())
	}

	return sb.String()
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

	extension := filepath.Ext(file.Filename)
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
