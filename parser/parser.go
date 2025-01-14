package parser

import (
	"0x822a5b87/test-fsm-arithmetic-operations"
	"0x822a5b87/test-fsm-arithmetic-operations/token"
	"fmt"
)

const (
	start fsm.State = iota
	end
	startGroup
	endGroup
	addOrSub
	mulOrDiv
	number
)

type ParserFsm struct {
	state  fsm.State
	acts   map[token.TokenType]fsm.Action
	tokens []token.Token
	index  int
	ast    *Ast
	// tmpAst temporary ast, used to save previous ast
	// there are 2 ways that makes us to build a temporary ast:
	// 1. encounter a lb;
	// 2. encounter an operator with lower precedence.
	tmpAst           *Ast
	prevOperatorType token.TokenType
}

func (p *ParserFsm) Parse() *Ast {
	for !p.isEnd() {
		curr := p.curr()
		p.Exec(curr)
	}
	return p.ast
}

func (p *ParserFsm) Exec(event token.Token) {
	action, ok := p.acts[event.Tt]
	if !ok {
		panic(fmt.Errorf("action not found for [%d]", event.Tt))
	}
	action(&event)
}

func (p *ParserFsm) AddAction(event token.Token, action fsm.Action) {
	p.acts[event.Tt] = action
}

func (p *ParserFsm) isEnd() bool {
	return p.curr().Tt == token.TokenEOF
}

func (p *ParserFsm) peek() token.Token {
	return p.tokens[p.index+1]
}

func (p *ParserFsm) curr() token.Token {
	return p.tokens[p.index]
}

func (p *ParserFsm) setState(state fsm.State) fsm.State {
	// TODO it should panic if an impossible state change occurs
	prevState := p.state
	p.state = state
	// TODO is this really correctly?
	p.index++
	return prevState
}

func (p *ParserFsm) inputNumber(event fsm.Event) {
	prevState := p.setState(number)
	numberToken := event.(*token.Token)
	if prevState == start {
		p.ast = NewAst(numberToken)
	} else if prevState == startGroup {
		// we've temporarily saved the previous ast in tmp ast, so it's a whole new ast to parse
		p.ast = NewAst(numberToken)
	} else {
		rhs := p.ast
		for rhs.Rhs != nil {
			rhs = rhs.Rhs
		}
		rhs.Rhs = NewAst(numberToken)
	}
}

func (p *ParserFsm) inputStartGroup(event fsm.Event) {
	p.prevOperatorType = event.(*token.Token).Tt
	p.setState(startGroup)
	p.tmpAst = p.ast
}

func (p *ParserFsm) inputEndGroup(event fsm.Event) {
	p.prevOperatorType = event.(*token.Token).Tt
	p.setState(endGroup)
	if p.tmpAst.Tk != nil {
		p.tmpAst.Rhs = p.ast
		p.ast = p.tmpAst
	}
}

func (p *ParserFsm) inputOperator(event fsm.Event) {
	p.prevOperatorType = event.(*token.Token).Tt
	operatorToken := event.(*token.Token)
	switch operatorToken.Tt {
	case token.TokenAdd:
		fallthrough
	case token.TokenSub:
		p.setState(addOrSub)
		ast := NewAst(operatorToken)
		ast.Lhs = p.ast
		p.ast = ast
	case token.TokenMul:
		fallthrough
	case token.TokenDiv:
		p.setState(mulOrDiv)
		ast := NewAst(operatorToken)
		if p.ast.Rhs == nil {
			ast.Lhs = p.ast
			p.ast = ast
		} else {
			ast.Lhs = p.ast.Rhs
			p.ast.Rhs = ast
		}
	default:
		panic(fmt.Errorf("unsupported operator type : [%d], value : [%s]", operatorToken.Tt, operatorToken.Value))
	}
}

func (p *ParserFsm) inputEnd(event fsm.Event) {
	p.setState(end)
}

func (p *ParserFsm) isEncounterLpOperator(event fsm.Event) bool {
	tk := event.(*token.Token)
	if tk.Tt == token.TokenAdd || tk.Tt == token.TokenSub {
		return p.prevOperatorType == token.TokenMul || p.prevOperatorType == token.TokenDiv
	}
	return false
}

func NewParserFsm(data string) *ParserFsm {
	tokenizerFsm := token.NewTokenizerFsm(data)
	tokens := tokenizerFsm.Tokenize()
	f := &ParserFsm{
		state:  start,
		acts:   make(map[token.TokenType]fsm.Action),
		tokens: tokens,
		index:  0,
		ast:    &Ast{},
	}

	f.acts[token.TokenNumber] = f.inputNumber
	f.acts[token.TokenLb] = f.inputStartGroup
	f.acts[token.TokenRb] = f.inputEndGroup
	f.acts[token.TokenAdd] = f.inputOperator
	f.acts[token.TokenSub] = f.inputOperator
	f.acts[token.TokenMul] = f.inputOperator
	f.acts[token.TokenDiv] = f.inputOperator
	f.acts[token.TokenEOF] = f.inputEnd

	return f
}
