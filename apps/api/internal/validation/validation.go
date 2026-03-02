package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Translate converts a go-playground/validator error into a human-readable
// Spanish message. Falls back to the raw error string for unknown tags.
func Translate(err error) string {
	var ve validator.ValidationErrors
	if !isValidationErrors(err, &ve) {
		return err.Error()
	}

	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		msgs = append(msgs, fieldMessage(fe))
	}
	return strings.Join(msgs, "; ")
}

func isValidationErrors(err error, ve *validator.ValidationErrors) bool {
	ok := false
	if e, ok2 := err.(validator.ValidationErrors); ok2 {
		*ve = e
		ok = true
	}
	return ok
}

func fieldMessage(fe validator.FieldError) string {
	field := spanishFieldName(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("El campo %s es obligatorio", field)
	case "email":
		return fmt.Sprintf("El campo %s debe ser un correo electrónico válido", field)
	case "min":
		if fe.Type().Kind().String() == "string" {
			return fmt.Sprintf("El campo %s debe tener al menos %s caracteres", field, fe.Param())
		}
		return fmt.Sprintf("El campo %s debe ser mayor o igual a %s", field, fe.Param())
	case "max":
		if fe.Type().Kind().String() == "string" {
			return fmt.Sprintf("El campo %s no puede superar los %s caracteres", field, fe.Param())
		}
		return fmt.Sprintf("El campo %s debe ser menor o igual a %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("El campo %s debe ser uno de: %s", field, strings.ReplaceAll(fe.Param(), " ", ", "))
	case "omitempty":
		return fmt.Sprintf("El campo %s tiene un valor inválido", field)
	default:
		return fmt.Sprintf("El campo %s no es válido (%s)", field, fe.Tag())
	}
}

// spanishFieldName maps Go struct field names to Spanish labels.
func spanishFieldName(f string) string {
	names := map[string]string{
		"Name":            "nombre",
		"Email":           "correo electrónico",
		"Password":        "contraseña",
		"Species":         "especie",
		"Breed":           "raza",
		"Age":             "edad",
		"Owner":           "dueño",
		"ConfirmPassword": "confirmación de contraseña",
	}
	if s, ok := names[f]; ok {
		return s
	}
	return strings.ToLower(f)
}
