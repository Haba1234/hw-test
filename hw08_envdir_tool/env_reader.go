package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envMap := make(Environment)

	for _, file := range files {
		// В имени файла не должно быть символа "=".
		if strings.Contains(file.Name(), "=") {
			continue
		}

		f, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		info, err := f.Stat()
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			// Пустые файлы помечаем к удалению.
			envMap[file.Name()] = EnvValue{"", true}
		} else {
			// Читаем только первую строку.
			scanner := bufio.NewScanner(f)
			if scanner.Scan() {
				str := scanner.Text()
				str = strings.TrimRight(str, " \t")
				str = string(bytes.Replace([]byte(str), []byte("\x00"), []byte("\n"), -1)) //nolint:gocritic
				envMap[file.Name()] = EnvValue{str, false}
			}
		}
		err = f.Close()
		if err != nil {
			return nil, err
		}
	}
	return envMap, nil
}
