package goutil

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validation struct {
	v     *validator.Validate
	trans *ut.Translator
}

func NewValidation() (Validation, error) {

	// translator
	en := en.New()
	uni := ut.New(en, en)

	// register translator
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	// send validate connection
	return Validation{validate, &trans}, nil

}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationErrors) IsErrorExists() bool {
	return len(v.Errors) > 0
}

func (vt *Validation) ValidationStruct(req interface{}) (validationError ValidationErrors) {

	if err := vt.v.Struct(req); err != nil {

		for _, errs := range err.(validator.ValidationErrors) {

			field, _ := reflect.TypeOf(req).FieldByName(errs.Field())
			filedJSONName, _ := field.Tag.Lookup("json")
			validationError.Errors = append(validationError.Errors, ValidationError{
				Field:   filedJSONName,
				Message: strings.Replace(errs.Translate(*vt.trans), errs.Field(), filedJSONName, -1),
			})
		}
	}

	return
}

func (vt *Validation) ValidationVariable(req interface{}, attributeName string, tag string, msgErr string) (validationError ValidationErrors) {

	if err := vt.v.Var(req, tag); err != nil {
		validationError.Errors = append(validationError.Errors, ValidationError{
			Field:   attributeName,
			Message: msgErr,
		})
	}

	return
}

type RegisterTranslation func(v *validator.Validate, trans *ut.Translator) error

func (vt *Validation) RegisterTranslation(translations ...RegisterTranslation) (err error) {

	registerJPG := func(ut ut.Translator) error { return ut.Add("jpg", "{0} must be a valid format", true) }
	translationJPG := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("jpg", fe.Field())
		return t
	}

	err = vt.v.RegisterTranslation("jpg", *vt.trans, registerJPG, translationJPG)
	if err != nil {
		return err
	}

	registerPNG := func(ut ut.Translator) error { return ut.Add("png", "{0} must be a valid format", true) }
	translationPNG := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("png", fe.Field())
		return t
	}

	err = vt.v.RegisterTranslation("png", *vt.trans, registerPNG, translationPNG)
	if err != nil {
		return err
	}

	registerPDF := func(ut ut.Translator) error { return ut.Add("pdf", "{0} must be a valid format", true) }
	translationPDF := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("pdf", fe.Field())
		return t
	}

	err = vt.v.RegisterTranslation("pdf", *vt.trans, registerPDF, translationPDF)
	if err != nil {
		return err
	}

	registerImage := func(ut ut.Translator) error { return ut.Add("image", "{0} must be a valid format", true) }
	translationImage := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("image", fe.Field())
		return t
	}

	err = vt.v.RegisterTranslation("image", *vt.trans, registerImage, translationImage)
	if err != nil {
		return err
	}

	for _, transFunc := range translations {
		err = transFunc(vt.v, vt.trans)
		if err != nil {
			return err
		}
	}

	return nil

}

type RegisterValidation func(v *validator.Validate) error

func (vt *Validation) RegisterValidation(validations ...RegisterValidation) (err error) {

	validatorJPG := func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^.(jpg|jpeg|JPG|JPEG)$`)
		return charValidation.MatchString(fl.Field().String())
	}

	err = vt.v.RegisterValidation("jpg", validatorJPG)
	if err != nil {
		return err
	}

	validatorPNG := func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^.(png|PNG)$`)
		return charValidation.MatchString(fl.Field().String())
	}

	err = vt.v.RegisterValidation("png", validatorPNG)
	if err != nil {
		return err
	}

	validatorPDF := func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^.(pdf|PDF)$`)
		return charValidation.MatchString(fl.Field().String())
	}

	err = vt.v.RegisterValidation("pdf", validatorPDF)
	if err != nil {
		return err
	}

	validatorImage := func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^.(jpg|jpeg|png|JPG|JPEG|PNG)$`)
		return charValidation.MatchString(fl.Field().String())
	}

	err = vt.v.RegisterValidation("image", validatorImage)
	if err != nil {
		return err
	}

	for _, validateFunc := range validations {
		err = validateFunc(vt.v)
		if err != nil {
			return err
		}
	}

	return nil

}
