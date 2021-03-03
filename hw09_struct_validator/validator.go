package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrValidateWrong = errors.New("validate wrong")
var ErrFieldWrong = errors.New("field wrong")

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	/*var errStr string
	for _, err := range v{
		errStr += fmt.Sprintf("Field: %s: %s\n", err.Field, err.Err)
	}
	return errStr*/

	var b []error
	for _, val := range v {
		b = append(b, fmt.Errorf("%s: %w\n", val.Field, val.Err))
	}

	return fmt.Sprintf("%v", b)
}

type fieldTag struct {
	fieldName string
	fieldVal  reflect.Value
	tagVal    reflect.StructTag
}

// Проверка строки по ключу 'len:xxx'.
func validateKeyLen(s string, key string) error {
	valKey, err := strconv.Atoi(regexp.MustCompile(`([0-9]+)`).FindString(key))
	if err != nil {
		return err
	}
	fmt.Printf("valKey: %v\n", valKey)
	if valKey < 0 {
		return ErrValidateWrong
	}
	if len(s) != valKey {
		return ErrFieldWrong
	}
	return nil
}

// Проверка на вхождение строки в множество строк. Ключ 'in:foo,bar'.
func validateKeyIn(value interface{}, key string) error {
	arrayStr := regexp.MustCompile(`(in:|,)`).Split(key, -1)
	arrayStr = arrayStr[1:]
	fmt.Println("arrayStr:", arrayStr)
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Int64 {
		for _, v := range arrayStr {
			tmp, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			if int64(tmp) == val.Int() {
				return nil
			}
		}
	}

	if val.Kind() == reflect.String {
		for _, v := range arrayStr {
			if v == val.String() {
				return nil
			}
		}
	}
	return ErrFieldWrong
}

// Проверка строки по ключу 'regexp:...'.
func validateKeyRegexp(s string, key string) error {
	re := regexp.MustCompile(`[^regexp:].*`).FindString(key)
	result, err := regexp.MatchString(re, s)
	fmt.Println(result, err)
	if err != nil {
		return err
	}
	if !result {
		return ErrFieldWrong
	}
	return nil
}

// Проверка числа на минимум, максимум Ключи 'min', 'max'.
// less = true - проверка на минимум.
func validateKeyMinMax(val int64, key string, less bool) error {
	valKey, err := strconv.Atoi(regexp.MustCompile(`([0-9]+)`).FindString(key))
	if err != nil {
		return err
	}
	fmt.Printf("valKey: %v\n", valKey)
	if less && val < int64(valKey) || !less && val > int64(valKey) {
		return ErrFieldWrong
	}
	return nil
}

// Парсим Tag и анализируем с типом поля
func (fl fieldTag) parsing(errValid *ValidationErrors) {
	var err error
	if keys, ok := fl.tagVal.Lookup("validate"); ok {
		if keys == "" {
			fmt.Println("(пустой тег)") // Добавить в ошибку
		} else {
			fmt.Println("Tag keys:", keys)
			fmt.Println(fl.fieldVal.Kind())
			// type string
			if fl.fieldVal.Kind() == reflect.String {
				// Выделение и проверка валидаторов.
				for _, key := range strings.Split(keys, "|") {
					switch {
					case strings.Contains(key, "len:"):
						err = validateKeyLen(fl.fieldVal.String(), key)
						fmt.Println("!!!result!!!", err)
					case strings.Contains(key, "regexp:"):
						err = validateKeyRegexp(fl.fieldVal.String(), key)
						fmt.Println("!!!result!!!", err)
					case strings.Contains(key, "in:"):
						err = validateKeyIn(fl.fieldVal.String(), key)
						fmt.Println("!!!result!!!", err)
					default:
						fmt.Println("Неправильный валидатор!!!")
					}
					fmt.Printf("error: %v\n", err)
					if err != nil {
						*errValid = append(*errValid, ValidationError{
							Field: fl.fieldName,
							Err:   err,
						})
					}
					fmt.Println("errValid:", errValid)
				}
			}
			// type int
			if fl.fieldVal.Kind() == reflect.Int {
				// Выделение и проверка валидаторов.
				for _, key := range strings.Split(keys, "|") {
					switch {
					case strings.Contains(key, "min:"):
						err = validateKeyMinMax(fl.fieldVal.Int(), key, true)
						fmt.Println("!!!result!!!", err)
					case strings.Contains(key, "max:"):
						err = validateKeyMinMax(fl.fieldVal.Int(), key, false)
						fmt.Println("!!!result!!!", err)
					case strings.Contains(key, "in:"):
						err = validateKeyIn(fl.fieldVal.Int(), key)
						fmt.Println("!!!result!!!", err)
					default:
						fmt.Println("Неправильный валидатор!!!")
					}
					if err != nil {
						*errValid = append(*errValid, ValidationError{
							Field: fl.fieldName,
							Err:   err,
						})
					}
					fmt.Println("errValid:", errValid)
				}
			}
		}
	}
	return
}

func Validate(v interface{}) error {
	var errValid ValidationErrors
	valStruct := reflect.ValueOf(v)
	fmt.Printf("type : %T\n", v)
	if valStruct.Kind() != reflect.Struct {
		fmt.Printf("expected a struct, but received %T\n", v)
		return nil
	}
	fmt.Println("type:", reflect.TypeOf(v))
	typeStruct := valStruct.Type()
	fmt.Println("NumField:", typeStruct.NumField()) // Кол-во полей
	for i := 0; i < typeStruct.NumField(); i++ {
		fmt.Println("reflect.StructField:", typeStruct.Field(i))
		fmt.Println("reflect.StructField.Name:", typeStruct.Field(i).Name)
		fmt.Println("reflect.StructField.Tag:", typeStruct.Field(i).Tag)
		fieldTag{
			fieldName: typeStruct.Field(i).Name,
			fieldVal:  valStruct.Field(i),
			tagVal:    typeStruct.Field(i).Tag,
		}.parsing(&errValid)
	}
	return errValid
}
