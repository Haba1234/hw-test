package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type progressBar struct {
	percent      int64  // Процент выполнения.
	cur          int64  // Текущее значение.
	total        int64  // Всего надо выполнить.
	rate         string // Текущая позиция символа.
	symbol       string // Символ заполнения.
	repeatSymbol int    // Кол-во повторов символа, если total < lenBar.
}

const lenBar = 50 // Длина в символах progressBar.

func newOption(start, total int64, symbol string) *progressBar {
	// Пересчет шкалы, если байт для копирования меньше lenBar (lenBar - длина progress Bar).
	repeatSymbol := 1
	if total < lenBar {
		repeatSymbol = int(math.Round(lenBar / float64(total)))
	}
	return &progressBar{
		cur:          start,
		total:        total,
		symbol:       symbol,
		repeatSymbol: repeatSymbol,
	}
}

func (bar *progressBar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *progressBar) increment(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += strings.Repeat(bar.symbol, bar.repeatSymbol)
	}

	fmt.Printf("\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent, bar.cur, bar.total) //nolint:forbidigo
}

func (bar *progressBar) finish(textMsg string, duration time.Duration) {
	fmt.Println()                                                                   //nolint:forbidigo
	fmt.Println(textMsg)                                                            //nolint:forbidigo
	fmt.Printf("Total: %d/%d bytes, copy time: %v\n", bar.cur, bar.total, duration) //nolint:forbidigo
}

//nolint:funlen,cyclop
func Copy(fromPath, toPath string, offset, limit int64) error {
	// Проверка на неверные ключи.
	if len(fromPath) == 0 {
		return fmt.Errorf("%w. Wrong key 'from': %q", ErrUnsupportedFile, fromPath)
	}
	if len(toPath) == 0 {
		return fmt.Errorf("%w. Wrong key 'to': %q", ErrUnsupportedFile, toPath)
	}

	f, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	sizeFile := info.Size()

	if sizeFile == 0 {
		return fmt.Errorf("wrong key 'from': %q, unknown file length. %w ", fromPath, ErrUnsupportedFile)
	}
	if offset > sizeFile {
		return fmt.Errorf("wrong key 'offset' = %d, file size = %d. %w ", offset, sizeFile, ErrOffsetExceedsFileSize)
	}
	if offset > 0 {
		_, err = f.Seek(offset, 0)
		if err != nil {
			return err
		}
	}
	if limit == 0 {
		limit = sizeFile
	}

	c, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer c.Close()

	bar := newOption(0, limit, "#")
	start := time.Now()
	for i := int64(1); i <= limit; i++ {
		_, err = io.CopyN(c, f, 1)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond)
		bar.increment(i)
	}
	duration := time.Since(start)
	text := "Copy finish!"

	if err != nil {
		if errors.Is(err, io.EOF) {
			text = "End of file reached!"
		} else {
			text = "Copy error!"
		}
	}

	bar.finish(text, duration)

	if !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
