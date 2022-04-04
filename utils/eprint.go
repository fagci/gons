package utils

import (
	"fmt"
	"os"
)

func EPrint(s ...interface{}) {
	os.Stderr.WriteString(fmt.Sprint(s...))
}

func EPrintf(f string, s ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(f, s...))
}

func EPrintln(s ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintln(s...))
}
