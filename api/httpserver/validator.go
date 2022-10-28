package httpserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator() (*Validator, error) {
	uni := ut.New(en.New())
	translation, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("error on get translador, 'en'not found")
	}
	validate := validator.New()

	err := en_translations.RegisterDefaultTranslations(validate, translation)
	if err != nil {
		return nil, err
	}
	return &Validator{validate, translation}, nil
}

func validateRequestBody(err error, w http.ResponseWriter, translator ut.Translator) {
	validationErrors := err.(validator.ValidationErrors)
	errorResponses := make([]ErrorResponse, len(validationErrors))
	for i, vErr := range validationErrors {
		errorResponses[i] = NewValidationError(vErr.Translate(translator))
	}
	apiResponse(w, http.StatusBadRequest, NewErrorResponses(errorResponses...))
}

func NewValidationError(message string) ErrorResponse {
	return ErrorResponse{Code: "validation", Message: strings.ToLower(message)}
}
