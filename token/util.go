package token

import "unicode"

func isDigit(event tokenizerEvent) bool {
	e := byte(event)
	return unicode.IsDigit(rune(e))
}
