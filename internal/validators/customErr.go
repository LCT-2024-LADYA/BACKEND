package validators

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
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
			case "password":
				message = fmt.Sprintf("Поле `%s` должно содержать от 8 до 64 символов и включать заглавные и строчные буквы, цифры и специальные символы.", jsonTag)
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
