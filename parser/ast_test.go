package parser

import (
	"0x822a5b87/test-fsm-arithmetic-operations/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	ast := Ast{
		Tk: &token.Token{
			Tt:    token.TokenAdd,
			Value: "+",
		},
		Lhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "123",
			},
		},

		Rhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "456",
			},
		},
	}
	assert.Equal(t, ast.Exec(), int64(579))
}

func TestSub(t *testing.T) {
	ast := Ast{
		Tk: &token.Token{
			Tt:    token.TokenSub,
			Value: "-",
		},
		Lhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "123",
			},
		},

		Rhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "456",
			},
		},
	}
	assert.Equal(t, ast.Exec(), int64(-333))
}

func TestMul(t *testing.T) {
	ast := Ast{
		Tk: &token.Token{
			Tt:    token.TokenMul,
			Value: "*",
		},
		Lhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "123",
			},
		},

		Rhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "456",
			},
		},
	}
	assert.Equal(t, ast.Exec(), int64(123*456))
}

func TestDiv(t *testing.T) {
	ast := Ast{
		Tk: &token.Token{
			Tt:    token.TokenDiv,
			Value: "/",
		},
		Lhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "789",
			},
		},

		Rhs: &Ast{
			Tk: &token.Token{
				Tt:    token.TokenNumber,
				Value: "263",
			},
		},
	}
	assert.Equal(t, ast.Exec(), int64(789/263))
}
