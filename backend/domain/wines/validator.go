package wines

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

type WineJSONValidator struct {
	v *validator.Validate
}

func NewWineJSONValidator() *WineJSONValidator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("vintage", ValidateVintage)
	return &WineJSONValidator{v}
}

func (v *WineJSONValidator) Validate(i interface{}) error {
	return v.v.Struct(i)
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
