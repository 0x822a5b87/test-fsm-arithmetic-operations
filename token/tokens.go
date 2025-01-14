package token

import "fmt"

const (
	TokenNumber TokenType = iota
	TokenAdd
	TokenSub
	TokenMul
	TokenDiv
	TokenLb
	TokenRb
	TokenEOF
)

type TokenType int

type Token struct {
	Tt    TokenType
	Value string
}

func (p *Token) EventType() int {
	return int(p.Tt)
}

func newOperatorToken(op string) Token {
	switch op {
	case "+":
		return Token{Tt: TokenAdd, Value: op}
	case "-":
		return Token{Tt: TokenSub, Value: op}
	case "*":
		return Token{Tt: TokenMul, Value: op}
	case "/":
		return Token{Tt: TokenDiv, Value: op}
	default:
		panic(fmt.Errorf("no such operator for [%s]", op))
	}
}

func newNumberToken(num string) Token {
	return Token{Tt: TokenNumber, Value: num}
}

func newLbToken() Token {
	return Token{
		Tt:    TokenLb,
		Value: "(",
	}
}

func newRbToken() Token {
	return Token{
		Tt:    TokenRb,
		Value: ")",
	}
}
