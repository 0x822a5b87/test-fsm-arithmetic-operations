package parser

import (
	"0x822a5b87/test-fsm-arithmetic-operations/token"
	"fmt"
	"strconv"
)

type Ast struct {
	Tk  *token.Token
	Lhs *Ast
	Rhs *Ast
}

func (ast *Ast) InsertToken(tk *token.Token) {

}

func (ast *Ast) Exec() int64 {
	switch ast.Tk.Tt {
	case token.TokenNumber:
		return ast.execNumber()
	case token.TokenAdd:
		return ast.execAdd()
	case token.TokenSub:
		return ast.execSub()
	case token.TokenMul:
		return ast.execMul()
	case token.TokenDiv:
		return ast.execDiv()
	default:
		panic(fmt.Errorf("unsupported type = [%d], value = [%s]", ast.Tk.Tt, ast.Tk.Value))
	}
}

func (ast *Ast) execNumber() int64 {
	num, _ := strconv.Atoi(ast.Tk.Value)
	return int64(num)
}

func (ast *Ast) execAdd() int64 {
	return ast.Lhs.Exec() + ast.Rhs.Exec()
}

func (ast *Ast) execSub() int64 {
	return ast.Lhs.Exec() - ast.Rhs.Exec()
}

func (ast *Ast) execMul() int64 {
	return ast.Lhs.Exec() * ast.Rhs.Exec()
}

func (ast *Ast) execDiv() int64 {
	return ast.Lhs.Exec() / ast.Rhs.Exec()
}

func NewAst(tk *token.Token) *Ast {
	return &Ast{
		Tk: tk,
	}
}
