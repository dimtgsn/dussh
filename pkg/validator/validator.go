package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
)

var validate *validator.Validate

func init() {
	Setup()
}

func Setup(options ...validator.Option) {
	validate = validator.New(options...)
}

func StructValidate(s any) error {
	var err error

	if kindOfData(s) != reflect.Struct {
		return nil
	}

	validateErrors := validate.Struct(s)
	if validateErrors != nil {
		if errs, ok := validateErrors.(validator.ValidationErrors); ok {
			for _, e := range errs {
				err = errors.Join(e)
			}
		}
	}

	return err
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
