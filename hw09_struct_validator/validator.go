package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	errorTemplate = "%w constraint %s"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrMin    = errors.New("minimum value error")
	ErrMax    = errors.New("maximum value error")
	ErrIn     = errors.New("contains error")
	ErrLen    = errors.New("string length error")
	ErrRegexp = errors.New("regexp error")
	ErrClient = errors.New("client error")
)

type Validator interface {
	IsValid(value reflect.Value) error
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, err := range v {
		builder.WriteString(fmt.Sprintf("%s: %s\n", err.Field, err.Error()))
	}
	return builder.String()
}

func (v ValidationErrors) Len() int {
	return len(v)
}

type MinValidator struct {
	Constraint int64
}

func (mv MinValidator) IsValid(value reflect.Value) error {
	if value.Kind() != reflect.Int {
		return fmt.Errorf("%w value is not int64", ErrClient)
	}
	if value.Int() < mv.Constraint {
		return ErrMin
	}
	return nil
}

type MaxValidator struct {
	Constraint int64
}

func (mv MaxValidator) IsValid(value reflect.Value) error {
	if value.Kind() != reflect.Int {
		return fmt.Errorf("%w value is not int64", ErrClient)
	}
	if value.Int() > mv.Constraint {
		return ErrMax
	}
	return nil
}

type LenValidator struct {
	Constraint int
}

func (lv LenValidator) IsValid(value reflect.Value) error {
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).Kind() != reflect.String {
				return fmt.Errorf("%w value is not string", ErrClient)
			}
			if value.Index(i).Len() != lv.Constraint {
				return ErrLen
			}
		}
	case reflect.String:
		if value.Len() != lv.Constraint {
			return ErrLen
		}
	default:
		return fmt.Errorf("%w value is not int", ErrClient)
	}
	return nil
}

type InValidator struct {
	Constraint []string
}

func (iv InValidator) IsValid(value reflect.Value) error {
	switch value.Kind() {
	case reflect.String:
		for _, v := range iv.Constraint {
			if value.String() == v {
				return nil
			}
		}
	case reflect.Int:
		for _, v := range iv.Constraint {
			item, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf(errorTemplate, ErrClient, v)
			}
			if value.Int() == int64(item) {
				return nil
			}
		}
	default:
		return fmt.Errorf("%w value is not string or int", ErrClient)
	}

	return ErrIn
}

type RegexpValidator struct {
	Constraint string
}

func (rv RegexpValidator) IsValid(value reflect.Value) error {
	if value.Kind() != reflect.String {
		return fmt.Errorf("%w value is not string", ErrClient)
	}
	matched, err := regexp.MatchString(rv.Constraint, value.String())
	if err != nil {
		return fmt.Errorf("%w value is not correct regexp", ErrClient)
	}
	if !matched {
		return ErrRegexp
	}
	return nil
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors
	validatingStruct := reflect.ValueOf(v)
	if validatingStruct.Kind() != reflect.Struct {
		return fmt.Errorf("%w value is not a struct", ErrClient)
	}
	for i := 0; i < validatingStruct.NumField(); i++ {
		typeField := validatingStruct.Type().Field(i)
		tag := typeField.Tag.Get("validate")
		if tag == "" {
			continue
		}
		validators, err := GetValidator(tag, typeField)
		if err != nil {
			return err
		}
		for _, validator := range validators {
			err := validator.IsValid(validatingStruct.Field(i))
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: typeField.Name,
					Err:   err,
				})
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func GetValidator(tag string, typeField reflect.StructField) ([]Validator, error) {
	validators := make([]Validator, 0)
	for _, constraint := range strings.Split(tag, "|") {
		validator, err := GetValidatorByConstraint(constraint, typeField)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validator)
	}
	return validators, nil
}

func GetValidatorByConstraint(constraint string, typeField reflect.StructField) (Validator, error) {
	split := strings.Split(constraint, ":")
	if len(split) < 2 {
		return nil, fmt.Errorf("%w invalid value provided for tag", ErrClient)
	}
	validator, ok := validators[split[0]]
	if !ok {
		return nil, fmt.Errorf(errorTemplate, ErrClient, constraint)
	}
	return validator(split[1], typeField)
}

var validators = map[string]func(string, reflect.StructField) (Validator, error){
	"min": func(constraint string, _ reflect.StructField) (Validator, error) {
		minimum, err := strconv.ParseInt(constraint, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(errorTemplate, ErrClient, constraint)
		}
		return MinValidator{Constraint: minimum}, nil
	},
	"max": func(constraint string, _ reflect.StructField) (Validator, error) {
		maximum, err := strconv.ParseInt(constraint, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(errorTemplate, ErrClient, constraint)
		}
		return MaxValidator{Constraint: maximum}, nil
	},
	"in": func(constraint string, _ reflect.StructField) (Validator, error) {
		values := strings.Split(constraint, ",")
		if len(values) < 2 {
			return nil, fmt.Errorf("%w invalid value provided for tag", ErrClient)
		}
		return InValidator{Constraint: values}, nil
	},
	"len": func(constraint string, _ reflect.StructField) (Validator, error) {
		length, err := strconv.Atoi(constraint)
		if err != nil {
			return nil, fmt.Errorf(errorTemplate, ErrClient, constraint)
		}
		return LenValidator{Constraint: length}, nil
	},
	"regexp": func(constraint string, _ reflect.StructField) (Validator, error) {
		_, err := regexp.Compile(constraint)
		if err != nil {
			return nil, fmt.Errorf(errorTemplate, ErrClient, constraint)
		}
		return RegexpValidator{Constraint: constraint}, nil
	},
}
