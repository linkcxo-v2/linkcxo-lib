package linkcxo

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	pgValidator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validate struct {
	pgValidate *pgValidator.Validate
	pgTrans    ut.Translator
}

type ValidationError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

func NewValidator() *Validate {
	v := Validate{
		pgValidate: pgValidator.New(),
	}
	var uni *ut.UniversalTranslator
	en := en.New()
	uni = ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.pgValidate, trans)
	v.pgTrans = trans
	return &v
}
func (v *Validate) Validate(val interface{}) []ValidationError {

	errors := []ValidationError{}
	err := v.pgValidate.Struct(val)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*pgValidator.InvalidValidationError); ok {
			fmt.Println(err)
			return errors
		} else {
			for _, err := range err.(pgValidator.ValidationErrors) {
				errors = append(errors, ValidationError{
					Field:   err.Field(),
					Tag:     err.Tag(),
					Value:   err.Value(),
					Message: err.Translate(v.pgTrans),
				})
			}
		}
	}
	return errors
}
