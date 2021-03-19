package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrValidateWrong = errors.New("validate wrong")
	ErrValidate      = errors.New("validated error")
	ErrNoStruct      = errors.New("invalid data type")
	reDigit          = regexp.MustCompile(`([0-9-]+)`)
	reIn             = regexp.MustCompile(`(in:|,)`)
	reRe             = regexp.MustCompile(`[^regexp:].*`)
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStr := strings.Builder{}
	for _, err := range v {
		// Оставляем только ошибки валидации
		if errors.Is(err.Err, ErrValidate) {
			errStr.WriteString(fmt.Sprintf("Field: %s: %s\n", err.Field, err.Err))
		}
	}
	return errStr.String()
}

type fieldT struct {
	fieldName string
	fieldVal  reflect.Value
	tagVal    reflect.StructTag
}

// Проверка строки по ключу 'len:xxx'.
func (field fieldT) validateKeyLen(key string) error {
	slice := reflect.ValueOf([]string{})
	// определяем тип данных: слайс или нет.
	switch field.fieldVal.Kind() {
	case reflect.String:
		slice = reflect.Append(slice, field.fieldVal)
	case reflect.Slice:
		slice = reflect.AppendSlice(slice, field.fieldVal.Slice(0, field.fieldVal.Len()))
	default:
		return ErrValidateWrong
	}

	valKey, err := strconv.Atoi(reDigit.FindString(key))
	if err != nil {
		return err
	}
	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		switch {
		case valKey < 0:
			return ErrValidateWrong
		case len(v.String()) != valKey:
			return ErrValidate
		}
	}
	return nil
}

// Проверка на вхождение строки в множество строк. Ключ 'in:foo,bar'.
func (field fieldT) validateKeyIn(key string) error {
	slice := reflect.ValueOf([]string{})
	// определяем тип данных: слайс или нет.
	switch field.fieldVal.Kind() {
	case reflect.String:
		slice = reflect.Append(slice, reflect.ValueOf(field.fieldVal.String()))
	case reflect.Int:
		slice = reflect.ValueOf([]int{})
		slice = reflect.Append(slice, field.fieldVal)
	case reflect.Slice:
		slice = reflect.AppendSlice(slice, field.fieldVal.Slice(0, field.fieldVal.Len()))
	default:
		return ErrValidateWrong
	}

	arrayKey := reIn.Split(key, -1)
	if len(arrayKey) < 2 {
		return ErrValidateWrong
	}
	arrayKey = arrayKey[1:]

	var valResult string
	// Проходим по слайсу значений.
	for i := 0; i < slice.Len(); i++ {
		val := slice.Index(i)

		switch val.Kind() {
		case reflect.String:
			valResult = val.String()
		case reflect.Int:
			valResult = strconv.Itoa(int(val.Int()))
		default:
			return ErrValidateWrong
		}
		// Сравниваем значение с ключом.
		keyOK := true
		for _, k := range arrayKey {
			if k == valResult {
				keyOK = false
				break
			}
		}
		if keyOK {
			return ErrValidate
		}
	}
	return nil
}

// Проверка строки по ключу 'regexp:...'.
func (field fieldT) validateKeyRegexp(key string) error {
	if field.fieldVal.Kind() == reflect.String {
		re := reRe.FindString(key)
		result, err := regexp.MatchString(re, field.fieldVal.String())
		if err != nil {
			return err
		}
		if !result {
			return ErrValidate
		}
	}
	return nil
}

// Проверка числа на минимум, максимум Ключи 'min', 'max'.
// less = true - проверка на минимум.
func (field fieldT) validateKeyMinMax(key string, less bool) error {
	if field.fieldVal.Kind() != reflect.Int {
		return ErrValidateWrong
	}
	valKey, err := strconv.Atoi(reDigit.FindString(key))
	if err != nil {
		return err
	}
	switch {
	case valKey < 0:
		return ErrValidateWrong
	case less && field.fieldVal.Int() < int64(valKey): // Проверка на min.
		return ErrValidate
	case !less && field.fieldVal.Int() > int64(valKey): // Проверка на max.
		return ErrValidate
	}
	return nil
}

// Парсим Tag и анализируем с типом поля.
func (field fieldT) parsing(errValid *ValidationErrors) {
	var err error

	if keys, ok := field.tagVal.Lookup("validate"); ok && keys != "" {
		// Выделение и проверка валидаторов.
		for _, key := range strings.Split(keys, "|") {
			switch {
			case strings.Contains(key, "len:"):
				err = field.validateKeyLen(key)
			case strings.Contains(key, "regexp:"):
				err = field.validateKeyRegexp(key)
			case strings.Contains(key, "min:"):
				err = field.validateKeyMinMax(key, true)
			case strings.Contains(key, "max:"):
				err = field.validateKeyMinMax(key, false)
			case strings.Contains(key, "in:"):
				err = field.validateKeyIn(key)
			default:
				err = nil
			}
			if err != nil {
				*errValid = append(*errValid, ValidationError{
					Field: field.fieldName,
					Err:   err,
				})
			}
		}
	}
}

func Validate(v interface{}) error {
	errValid := ValidationErrors{}
	valStruct := reflect.ValueOf(v)
	if valStruct.Kind() != reflect.Struct {
		return ErrNoStruct
	}

	typeStruct := valStruct.Type()
	for i := 0; i < typeStruct.NumField(); i++ {
		fieldT{
			fieldName: typeStruct.Field(i).Name,
			fieldVal:  valStruct.Field(i),
			tagVal:    typeStruct.Field(i).Tag,
		}.parsing(&errValid)
	}

	if len(errValid) != 0 {
		return errValid
	}
	return nil
}
