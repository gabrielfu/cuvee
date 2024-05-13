package wines

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

type WineJSONValidator struct {
	v *validator.Validate
}

type WineValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func tagToMessage(tag string) string {
	switch tag {
	case "required":
		return "Field is required"
	case "vintage":
		return "Vintage must be 'NV' or a year number"
	default:
		return ""
	}
}

func NewWineJSONValidator() *WineJSONValidator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("vintage", ValidateVintage)
	return &WineJSONValidator{v}
}

func (v *WineJSONValidator) Validate(i interface{}) []WineValidationError {
	err := v.v.Struct(i)
	if err != nil {
		ve := err.(validator.ValidationErrors)
		errs := make([]WineValidationError, len(ve))
		for i, err := range ve {
			errs[i] = WineValidationError{
				Field:   err.Field(),
				Message: tagToMessage(err.Tag()),
			}
		}
		return errs
	}
	return []WineValidationError{}
}

// ValidateVintage validates the vintage is either "NV" or an integer
func ValidateVintage(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	if value == "NV" {
		return true
	}
	_, err := strconv.Atoi(value)
	return err == nil
}
