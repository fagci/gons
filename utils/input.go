package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/fagci/gons/loaders"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func LoadInput(input string) ([]string, error) {
	var lines []string
	var err error

	switch {
	case input == "-":
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			lines = append(lines, reader.Text())
		}
		err = reader.Err()
	case fileExists(input):
		lines, err = loaders.FileToArray(input)
	default:
		lines = strings.FieldsFunc(input, func(c rune) bool {
			return c == '\n' || c == '\r'
		})
	}

	return lines, err
}
