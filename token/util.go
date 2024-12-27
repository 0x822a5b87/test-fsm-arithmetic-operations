package token

import "unicode"

func isDigit(event TokenizerEvent) bool {
	e := byte(event)
	return unicode.IsDigit(rune(e))
}
