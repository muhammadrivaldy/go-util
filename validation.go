package goutil

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gobeam/stringy"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type validation struct {
	v     *validator.Validate
	trans ut.Translator
}

func NewValidation() (validation, error) {

	// translator
	en := en.New()
	uni := ut.New(en, en)

	// register translator
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	// send validate connection
	return validation{validate, trans}, nil

}

func (vt *validation) ValidationStruct(req interface{}) error {

	if err := vt.v.Struct(req); err != nil {

		var temp []string

		for _, errs := range err.(validator.ValidationErrors) {
			newField := stringy.New(errs.Field()).SnakeCase().ToLower()
			tempMessage := strings.Replace(errs.Translate(vt.trans), errs.Field(), newField, -1)
			temp = append(temp, tempMessage)
		}

		return errors.New(strings.Join(temp, ", "))

	}

	return nil

}

func (vt *validation) ValidationVariable(req interface{}, tag string, msgErr string) error {

	if err := vt.v.Var(req, tag); err != nil {
		return errors.New(msgErr)
	}

	return nil

}

type RegisterTranslation func(v *validator.Validate, trans ut.Translator) error

func (vt *validation) RegisterTranslation(translations ...RegisterTranslation) (err error) {

	registerJPG := func(ut ut.Translator) error { return ut.Add("jpg", "{0} must be a valid format", true) }
	translationJPG := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("jpg", fe.Field())
		return t
	}

	err = vt.v.RegisterTranslation("jpg", vt.trans, registerJPG, translationJPG)
	if err != nil {
		return err
	}

	registerPNG := func(ut ut.Translator) error { return ut.Add("png", "{0} must be a valid format", true) }
	translationPNG := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("png", fe.Field())
		return t
	}

	err = vt.v.RegisterTranslation("png", vt.trans, registerPNG, translationPNG)
	if err != nil {
		return err
	}

	registerPDF := func(ut ut.Translator) error { return ut.Add("pdf", "{0} must be a valid format", true) }
	translationPDF := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("pdf", fe.Field())
		return t
	}

	err = vt.v.RegisterTranslation("pdf", vt.trans, registerPDF, translationPDF)
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

func (vt *validation) RegisterValidation(validations ...RegisterValidation) (err error) {

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

	for _, validateFunc := range validations {
		err = validateFunc(vt.v)
		if err != nil {
			return err
		}
	}

	return nil

}
