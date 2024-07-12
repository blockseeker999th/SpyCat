package validation

import "github.com/go-playground/validator/v10"

func ValidateStruct(validationStruct interface{}) error {
	err := validator.New().Struct(validationStruct)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}
