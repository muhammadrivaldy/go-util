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

// NewValidation ..
func NewValidation() (Validation, error) {

	// translator
	en := en.New()
	uni := ut.New(en, en)

	// register translator
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	// register validation
	if err := registerValidation(validate); err != nil {
		return &validation{}, err
	}

	// register translation
	if err := registerTranslation(validate, trans); err != nil {
		return &validation{}, err
	}

	// send validate connection
	return &validation{
		v:     validate,
		trans: trans}, nil

}

// ValidationStruct ..
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

// ValidationVariable ..
func (vt *validation) ValidationVariable(req interface{}, tag string, msgErr string) error {
	if err := vt.v.Var(req, tag); err != nil {
		return errors.New(msgErr)
	}
	return nil
}

// Validation ..
type Validation interface {
	ValidationStruct(req interface{}) error
	ValidationVariable(req interface{}, tag string, msgErr string) error
}

// register translation
func registerTranslation(validate *validator.Validate, trans ut.Translator) error {

	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("number", trans, func(ut ut.Translator) error {
		return ut.Add("number", "{0} must be number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("number", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} must be a valid phone", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("jpg", trans, func(ut ut.Translator) error {
		return ut.Add("jpg", "{0} must be a valid format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("jpg", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("png", trans, func(ut ut.Translator) error {
		return ut.Add("png", "{0} must be a valid format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("png", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("pdf", trans, func(ut ut.Translator) error {
		return ut.Add("pdf", "{0} must be a valid format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("pdf", fe.Field())
		return t
	})

	return nil

}

// register validation
func registerValidation(validate *validator.Validate) error {

	if err := validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^0(\d{0,12})$`)
		return charValidation.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	if err := validate.RegisterValidation("jpg", func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^.(jpg|jpeg|JPG|JPEG)$`)
		return charValidation.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	if err := validate.RegisterValidation("png", func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^.(png|PNG)$`)
		return charValidation.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	if err := validate.RegisterValidation("pdf", func(fl validator.FieldLevel) bool {
		charValidation := regexp.MustCompile(`^.(pdf|PDF)$`)
		return charValidation.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	return nil

}
