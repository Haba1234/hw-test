package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func multiChar(char rune, count int) string {
	return strings.Repeat(string(char), count)
}

func checkSymbols(char rune) bool { //nolint:gofmt
	// проверка, что символ соответствует разрешенному.
	result := true
	if !unicode.IsLetter(char) && !unicode.IsDigit(char) && unicode.IsPunct(char) && (string(char) != `\`) {
		result = false
	}
	return result
}

func Unpack(str string) (string, error) { //nolint:gocognit,funlen,gocyclo
	var prevChar rune
	var count int
	var punct, double bool
	var resultStr strings.Builder // готовим построитель для новой строки

	resultStr.Reset()

	// если строка пустая, сразу на выход
	if len(str) == 0 {
		return "", nil
	}

	// перебор символов строки
	for index, char := range str {
		// проверка на неразрешенные символы
		if !checkSymbols(char) {
			return "", ErrInvalidString
		}
		// строка начинается с цифры или \ - ошибка. Остальные символы отсеиваются в функции checkSymbols
		if index == 0 && (unicode.IsDigit(char) || unicode.IsPunct(char)) {
			return "", ErrInvalidString
		}

		if index > 0 {
			switch {
			case unicode.IsLetter(char) && unicode.IsLetter(prevChar):
				// буква после буквы
				resultStr.WriteString(multiChar(prevChar, 1))
			case unicode.IsPunct(char) && unicode.IsLetter(prevChar):
				// \ после буквы
				resultStr.WriteString(multiChar(prevChar, 1))
				punct = true
			case unicode.IsPunct(char) && unicode.IsDigit(prevChar) && punct:
				// \ после цифры
				resultStr.WriteString(multiChar(prevChar, 1))
				punct = true
			case unicode.IsPunct(char) && unicode.IsPunct(prevChar):
				// \ после \
				if double {
					resultStr.WriteString(multiChar(prevChar, 1))
					double = false
				} else {
					double = true
				}
			case unicode.IsLetter(char) && unicode.IsDigit(prevChar):
			case (unicode.IsDigit(char) && unicode.IsLetter(prevChar)) || (unicode.IsDigit(char) && unicode.IsDigit(prevChar) && punct) || double:
				// цифра после буквы, размножение символа
				count, _ = strconv.Atoi(string(char))
				resultStr.WriteString(multiChar(prevChar, count))
				punct = false
			case unicode.IsDigit(char) && unicode.IsPunct(prevChar):
				// цифра после \
			default: // в любых других случаях - ошибка
				return "", ErrInvalidString
			}
		}
		prevChar = char
	}

	if unicode.IsLetter(prevChar) || punct {
		// сохранить последний символ в строке перед возвратом результата
		resultStr.WriteString(multiChar(prevChar, 1))
	}

	return resultStr.String(), nil
}
