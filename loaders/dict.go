package loaders

import (
	"bufio"
	"os"
)

func FileToArray(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var items []string
	for scanner.Scan() {
		items = append(items, scanner.Text())
	}
	return items, scanner.Err()
}
