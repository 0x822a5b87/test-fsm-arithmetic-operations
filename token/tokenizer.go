package token

import (
	"0x822a5b87/test-fsm-arithmetic-operations"
	"bytes"
	"fmt"
)

const (
	startGroup  fsm.State = 1
	endGroup    fsm.State = 2
	newOperator fsm.State = 3
	appendDigit fsm.State = 4
	end         fsm.State = 5

	null    tokenizerEvent = 0
	space   tokenizerEvent = ' '
	tab     tokenizerEvent = '\t'
	r       tokenizerEvent = '\r'
	newLine tokenizerEvent = '\n'
	lb      tokenizerEvent = '('
	rb      tokenizerEvent = ')'
	add     tokenizerEvent = '+'
	sub     tokenizerEvent = '-'
	mul     tokenizerEvent = '*'
	div     tokenizerEvent = '/'
)

type tokenizerEvent byte

func (t tokenizerEvent) EventType() int {
	return int(t)
}

type TokenizerFsm struct {
	state  fsm.State
	acts   map[tokenizerEvent]fsm.Action
	stream charStream
	tokens []Token
}

func (fsm *TokenizerFsm) Tokenize() []Token {
	for !fsm.isEnd() {
		event := fsm.stream.peekEvent()
		fsm.Exec(event)
	}
	fsm.tokens = append(fsm.tokens, Token{Tt: TokenEOF, Value: ""})
	return fsm.tokens
}

func (fsm *TokenizerFsm) Exec(event tokenizerEvent) {
	action, ok := fsm.acts[event]
	if !ok {
		panic(fmt.Errorf("action not found for [%d]", event))
	}
	action(event)
}

func (fsm *TokenizerFsm) AddAction(event tokenizerEvent, action fsm.Action) {
	fsm.acts[event] = action
}

func (fsm *TokenizerFsm) setState(newState fsm.State) {
	fsm.state = newState
}

func (fsm *TokenizerFsm) skip(event fsm.Event) {
	fsm.stream.nextEvent()
}

func (fsm *TokenizerFsm) startGroup(event fsm.Event) {
	fsm.setState(startGroup)
	event = event.(tokenizerEvent)
	fsm.stream.nextEvent()
	fsm.tokens = append(fsm.tokens, newLbToken())
}

func (fsm *TokenizerFsm) endGroup(event fsm.Event) {
	fsm.setState(endGroup)
	event = event.(tokenizerEvent)
	fsm.stream.nextEvent()
	fsm.tokens = append(fsm.tokens, newRbToken())
}

func (fsm *TokenizerFsm) isEnd() bool {
	return fsm.state == end
}

func (fsm *TokenizerFsm) end(event fsm.Event) {
	fsm.setState(end)
}

func (fsm *TokenizerFsm) operator(event fsm.Event) {
	fsm.setState(newOperator)
	fsm.stream.nextEvent()

	e := event.(tokenizerEvent)
	buff := bytes.NewBufferString("")
	buff.WriteByte(byte(e))
	fsm.tokens = append(fsm.tokens, newOperatorToken(buff.String()))
}

func (fsm *TokenizerFsm) appendDigit(event fsm.Event) {
	fsm.setState(appendDigit)
	buff := bytes.NewBufferString("")
	for {
		e := fsm.stream.peekEvent()
		if !isDigit(e) {
			break
		}

		e = fsm.stream.nextEvent()
		buff.WriteByte(byte(e))
	}

	fsm.tokens = append(fsm.tokens, newNumberToken(buff.String()))
}

func NewTokenizerFsm(data string) *TokenizerFsm {
	tf := &TokenizerFsm{
		state:  startGroup,
		acts:   make(map[tokenizerEvent]fsm.Action),
		stream: charStream{data: data},
		tokens: make([]Token, 0),
	}

	tf.AddAction(null, tf.end)

	tf.AddAction(space, tf.skip)
	tf.AddAction(tab, tf.skip)
	tf.AddAction(r, tf.skip)
	tf.AddAction(newLine, tf.skip)

	tf.AddAction(lb, tf.startGroup)

	tf.AddAction(rb, tf.endGroup)

	tf.AddAction(add, tf.operator)
	tf.AddAction(sub, tf.operator)
	tf.AddAction(mul, tf.operator)
	tf.AddAction(div, tf.operator)

	chars := "0123456789"
	for _, ch := range chars {
		tf.AddAction(tokenizerEvent(ch), tf.appendDigit)
	}

	return tf
}
