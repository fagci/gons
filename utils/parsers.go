package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseRange(r string) []int {
	var values []int
	if len(r) == 0 {
		return values
	}
	parts := strings.Split(r, ",")
	for _, part := range parts {
		fromTo := strings.Split(part, "-")
		from, errf := strconv.Atoi(fromTo[0])
		if errf != nil {
			panic(fmt.Sprintf("Bad range:%s", part))
		}
		if len(fromTo) == 2 {
			to, errt := strconv.Atoi(fromTo[1])
			if errt != nil {
				panic(fmt.Sprintf("Bad range:%s", part))
			}

			for i := from; i <= to; i++ {
				values = append(values, i)
			}
			continue
		}

		values = append(values, from)
	}
	return values
}
