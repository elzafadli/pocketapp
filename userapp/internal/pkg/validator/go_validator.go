package validator

import (
	"context"
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type GoValidator struct {
	validate *validator.Validate
	uni      ut.Translator
}

func NewGoValidator() Validator {
	v := validator.New()
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(v, trans)

	return &GoValidator{validate: v, uni: trans}
}

func (v *GoValidator) Validate(ctx context.Context, data any) error {
	err := v.validate.StructCtx(ctx, data)
	if err == nil {
		return nil
	}

	// Check for invalid validation error (e.g., nil pointer, wrong type)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	// Safely assert to ValidationErrors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// If it's not ValidationErrors, return the error as-is
		return err
	}

	if len(errs) > 0 {
		mapErr := make(map[string]error, len(errs))
		for _, err := range errs {
			mapErr[err.Field()] = errors.New(err.Translate(v.uni))
		}

		return NewErrorMap(mapErr)
	}

	return nil
}
