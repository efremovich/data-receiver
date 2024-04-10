package aconf

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// validateStruct валидирует поля переданной структуры, у которых
// есть тэг "validate"
func validateStruct(v interface{}) error {
	err := validator.New().Struct(v)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrValidate, err)
	}
	return err
}
