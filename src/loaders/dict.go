package loaders

import (
	"bufio"
	"os"
)

type DictLoader struct {
	items []string
}

func NewDictLoader() *DictLoader {
	return &DictLoader{}
}

func (dictLoader *DictLoader) Load(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictLoader.items = append(dictLoader.items, scanner.Text())
	}
	return dictLoader.items, scanner.Err()
}
