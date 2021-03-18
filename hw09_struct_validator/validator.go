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
	reDigit   *regexp.Regexp
	reIn      *regexp.Regexp
	reRe      *regexp.Regexp
}

// Добавляет найденные ошибки в слайс.
func (field fieldT) appendValidateErr(errValid *ValidationErrors, err error) {
	if err != nil {
		*errValid = append(*errValid, ValidationError{
			Field: field.fieldName,
			Err:   err,
		})
	}
}

// Проверка строки по ключу 'len:xxx'.
func (field fieldT) validateKeyLen(errValid *ValidationErrors, key string) {
	slice := reflect.ValueOf([]string{})
	// определяем тип данных: слайс или нет.
	switch field.fieldVal.Kind() {
	case reflect.String:
		slice = reflect.Append(slice, field.fieldVal)
	case reflect.Slice:
		slice = reflect.AppendSlice(slice, field.fieldVal.Slice(0, field.fieldVal.Len()))
	default:
		field.appendValidateErr(errValid, ErrValidateWrong)
		return
	}

	valKey, err := strconv.Atoi(field.reDigit.FindString(key))
	if err != nil {
		field.appendValidateErr(errValid, err)
		return
	}
	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		switch {
		case valKey < 0:
			field.appendValidateErr(errValid, ErrValidateWrong)
		case len(v.String()) != valKey:
			field.appendValidateErr(errValid, ErrValidate)
		}
	}
}

// Проверка на вхождение строки в множество строк. Ключ 'in:foo,bar'.
func (field fieldT) validateKeyIn(errValid *ValidationErrors, key string) {
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
		field.appendValidateErr(errValid, ErrValidateWrong)
		return
	}

	arrayKey := field.reIn.Split(key, -1)
	if len(arrayKey) < 2 {
		field.appendValidateErr(errValid, ErrValidateWrong)
		return
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
			field.appendValidateErr(errValid, ErrValidateWrong)
			return
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
			field.appendValidateErr(errValid, ErrValidate)
		}
	}
}

// Проверка строки по ключу 'regexp:...'.
func (field fieldT) validateKeyRegexp(key string) error {
	if field.fieldVal.Kind() == reflect.String {
		re := field.reRe.FindString(key)
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
	valKey, err := strconv.Atoi(field.reDigit.FindString(key))
	if err != nil {
		return err
	}
	if less && field.fieldVal.Int() < int64(valKey) || !less && field.fieldVal.Int() > int64(valKey) {
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
				field.validateKeyLen(errValid, key)
			case strings.Contains(key, "regexp:"):
				err = field.validateKeyRegexp(key)
			case strings.Contains(key, "min:"):
				err = field.validateKeyMinMax(key, true)
			case strings.Contains(key, "max:"):
				err = field.validateKeyMinMax(key, false)
			case strings.Contains(key, "in:"):
				field.validateKeyIn(errValid, key)
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
		return errValid
	}

	typeStruct := valStruct.Type()
	for i := 0; i < typeStruct.NumField(); i++ {
		fieldT{
			fieldName: typeStruct.Field(i).Name,
			fieldVal:  valStruct.Field(i),
			tagVal:    typeStruct.Field(i).Tag,
			reDigit:   regexp.MustCompile(`([0-9]+)`),
			reIn:      regexp.MustCompile(`(in:|,)`),
			reRe:      regexp.MustCompile(`[^regexp:].*`),
		}.parsing(&errValid)
	}

	if len(errValid) != 0 {
		return errValid
	}
	return nil
}
