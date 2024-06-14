package http

import "unicode"

// Write a function that transforms an input string into upper snake case.
// The input string will only contain upper camel-case letters.
// The output string should only contain upper snake-case letters.
// The output string should contain any underscores if the input has humps.
// The output string should not contain any underscores if the input has no humps.
// The output string should not contain any underscores if the input is only one letter.
// The output string should not contain any underscores if the input is empty.
func upperSnakeCaseFrom(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToUpper(r))
	}
	return string(result)
}
