package validation

import (
	"fmt"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

type ErrorTranslator struct {
	ENTranslator ut.Translator
	IDTranslator ut.Translator
}

func NewDefaultTranslator(v *validator.Validate) *ErrorTranslator {
	english := en.New()
	uni := ut.New(english, english)
	enTrans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(v, enTrans)

	id := id.New()
	uni = ut.New(id, id)
	idTrans, _ := uni.GetTranslator("id")
	_ = idTranslations.RegisterDefaultTranslations(v, idTrans)

	return &ErrorTranslator{
		ENTranslator: enTrans,
		IDTranslator: idTrans,
	}
}

func translateError(err error, trans ut.Translator) (errs []string) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Sprint(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}
