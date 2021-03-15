package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("Не верно заданы ключи") //nolint:forbidigo
		os.Exit(1)
	}

	dirEnv := os.Args[1]
	mapEnv, err := ReadDir(dirEnv)
	if err != nil {
		fmt.Printf("Ошибка чтения каталога %s. Error: %v\n", dirEnv, err) //nolint:forbidigo
		os.Exit(1)
	}
	code := RunCmd(os.Args[2:], mapEnv)
	os.Exit(code)
}
