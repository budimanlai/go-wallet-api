package library

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	en := en.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")

	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
}

func Validator() *validator.Validate {
	return validate
}

func ParseError(err error) []string {
	var result []string
	errs := err.(validator.ValidationErrors)
	emsg := errs.Translate(trans)
	for _, e := range emsg {
		result = append(result, e)
	}
	return result
}

func FirstError(err error) string {
	e := ParseError(err)
	return e[0]
}
