package validation

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name     string `validate:"required"`
	Email    string `validate:"email"`
	Age      int    `validate:"min=18"`
	Password string `validate:"min=8,max=20"`
	Species  string `validate:"oneof=dog cat bird"`
}

func TestTranslate_NonValidationError(t *testing.T) {
	result := Translate(errors.New("generic error message"))
	require.Contains(t, result, "generic error message")
}

func TestTranslate_Required(t *testing.T) {
	v := validator.New()
	obj := testStruct{}
	err := v.Struct(obj)
	require.Error(t, err)

	result := Translate(err)
	require.Contains(t, result, "El campo nombre es obligatorio")
}

func TestTranslate_Email(t *testing.T) {
	v := validator.New()
	obj := testStruct{Name: "Test", Email: "invalid", Age: 20, Password: "12345678", Species: "dog"}
	err := v.Struct(obj)
	require.Error(t, err)

	result := Translate(err)
	require.Contains(t, result, "correo electrónico")
	require.Contains(t, result, "válido")
}

func TestTranslate_Min(t *testing.T) {
	v := validator.New()
	obj := testStruct{Name: "Test", Email: "a@b.com", Age: 15, Password: "12345678", Species: "dog"}
	err := v.Struct(obj)
	require.Error(t, err)

	result := Translate(err)
	require.Contains(t, result, "edad")
	require.Contains(t, result, "mayor o igual a 18")
}

func TestTranslate_MinMaxString(t *testing.T) {
	v := validator.New()
	obj := testStruct{Name: "Test", Email: "a@b.com", Age: 20, Password: "short", Species: "dog"}
	err := v.Struct(obj)
	require.Error(t, err)

	result := Translate(err)
	require.Contains(t, result, "contraseña")
	require.Contains(t, result, "al menos 8 caracteres")
}

func TestTranslate_MaxString(t *testing.T) {
	v := validator.New()
	longPwd := strings.Repeat("a", 30)
	obj := testStruct{Name: "Test", Email: "a@b.com", Age: 20, Password: longPwd, Species: "dog"}
	err := v.Struct(obj)
	require.Error(t, err)

	result := Translate(err)
	require.Contains(t, result, "contraseña")
	require.Contains(t, result, "no puede superar los 20 caracteres")
}

func TestTranslate_OneOf(t *testing.T) {
	v := validator.New()
	obj := testStruct{Name: "Test", Email: "a@b.com", Age: 20, Password: "12345678", Species: "dragon"}
	err := v.Struct(obj)
	require.Error(t, err)

	result := Translate(err)
	require.Contains(t, result, "especie")
	require.Contains(t, result, "uno de: dog, cat, bird")
}

func TestTranslate_MultipleErrors(t *testing.T) {
	v := validator.New()
	obj := testStruct{}
	err := v.Struct(obj)

	result := Translate(err)
	lines := strings.Split(result, "; ")
	require.GreaterOrEqual(t, len(lines), 4)
}

func TestTranslate_FallbackFieldName(t *testing.T) {
	v := validator.New()
	type custom struct {
		SomeField string `validate:"required"`
	}
	err := v.Struct(custom{})
	require.Error(t, err)

	result := Translate(err)
	require.Contains(t, result, "somefield")
}
