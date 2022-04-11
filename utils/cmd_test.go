package utils

import "testing"

func TestFilterInQuotes(t *testing.T) {
	input := `"can't escape $tring"` + "`!`"
	expectedOutput := `cant escape tring!`

	output := FilterValueInQuotes(input)
	if output != expectedOutput {
		t.Errorf("Output not expected: `%s` != `%s`", output, expectedOutput)
	}
}
