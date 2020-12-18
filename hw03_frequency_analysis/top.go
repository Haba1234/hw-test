package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	// проверка на пустой текст
	if len(text) == 0 {
		return nil
	}

	reg := regexp.MustCompile(`[^-А-Яа-я ]`)
	newText := reg.ReplaceAllString(text, " ") // оставляем в тексте только буквы, цифры, "пробел" и "тире"

	sl := strings.Split(newText, ` `) // преобразуем в слайс разделяя слова по симовлу "пробел"
	textFilters := map[string]int{}   // готовим словарь для заполнения

	for _, str := range sl {
		if (len(str) > 0) && (str != "-") { // отсеиваем пустые строки и если строка = "-"
			str := strings.ToLower(str)            // делаем все символы маленькими
			if value, ok := textFilters[str]; ok { // если уже есть такой ключ, то +1 к значению
				textFilters[str] = value + 1
			} else {
				textFilters[str] = 1 // ключа еще не было, создаем новую запись в словарь
			}
		}
	}

	type tmpStruct struct {
		key string
		val int
	}
	sortStr := make([]tmpStruct, 0, len(textFilters)) // создаем слайс сразу нужного размера

	// заполняем слайс значениями для последующей сортировки по убыванию
	for key, value := range textFilters {
		sortStr = append(sortStr, tmpStruct{key, value})
	}

	// сортируем
	sort.Slice(sortStr, func(i, j int) bool {
		return sortStr[i].val > sortStr[j].val
	})

	// готовим слайс строк для вывода результата и заполняем его 10 наиболее часто встречающимися словами
	countTop := 10
	resultStrTop10 := make([]string, countTop)
	if len(sortStr) > countTop {
		for i := 0; i < countTop; i++ {
			resultStrTop10[i] = sortStr[i].key
		}
	} else {
		return nil
	}

	return resultStrTop10
}
