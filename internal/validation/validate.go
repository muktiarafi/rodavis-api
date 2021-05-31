package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/muktiarafi/rodavis-api/internal/api"
)

type Validator struct {
	Validator  *validator.Validate
	Translator *ErrorTranslator
}

func NewValidator(validator *validator.Validate, trans *ErrorTranslator) *Validator {
	return &Validator{
		Validator:  validator,
		Translator: trans,
	}
}

func (v *Validator) Validate(op string, i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {

		return api.NewValidationException(
			op,
			translateError(err, v.Translator.ENTranslator),
		)
	}
	return nil
}
