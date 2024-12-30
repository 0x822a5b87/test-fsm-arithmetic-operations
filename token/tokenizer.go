package token

import (
	"0x822a5b87/test-fsm-arithmetic-operations"
	"bytes"
	"fmt"
)

const (
	StartGroup  fsm.State = 1
	EndGroup    fsm.State = 2
	NewOperator fsm.State = 3
	AppendDigit fsm.State = 4
	End         fsm.State = 5

	Null    tokenizerEvent = 0
	Space   tokenizerEvent = ' '
	Tab     tokenizerEvent = '\t'
	Return  tokenizerEvent = '\r'
	NewLine tokenizerEvent = '\n'
	LB      tokenizerEvent = '('
	RB      tokenizerEvent = ')'
	Add     tokenizerEvent = '+'
	Sub     tokenizerEvent = '-'
	Mul     tokenizerEvent = '*'
	Div     tokenizerEvent = '/'
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
	fsm.setState(StartGroup)
	event = event.(tokenizerEvent)
	fsm.stream.nextEvent()
	fsm.tokens = append(fsm.tokens, newLbToken())
}

func (fsm *TokenizerFsm) endGroup(event fsm.Event) {
	fsm.setState(EndGroup)
	event = event.(tokenizerEvent)
	fsm.stream.nextEvent()
	fsm.tokens = append(fsm.tokens, newRbToken())
}

func (fsm *TokenizerFsm) end(event fsm.Event) {
	fsm.setState(End)
}

func (fsm *TokenizerFsm) operator(event fsm.Event) {
	fsm.setState(NewOperator)
	fsm.stream.nextEvent()

	e := event.(tokenizerEvent)
	buff := bytes.NewBufferString("")
	buff.WriteByte(byte(e))
	fsm.tokens = append(fsm.tokens, newOperatorToken(buff.String()))
}

func (fsm *TokenizerFsm) appendDigit(event fsm.Event) {
	fsm.setState(AppendDigit)
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
		state:  StartGroup,
		acts:   make(map[tokenizerEvent]fsm.Action),
		stream: charStream{data: data},
		tokens: make([]Token, 0),
	}

	tf.AddAction(Null, tf.end)

	tf.AddAction(Space, tf.skip)
	tf.AddAction(Tab, tf.skip)
	tf.AddAction(Return, tf.skip)
	tf.AddAction(NewLine, tf.skip)

	tf.AddAction(LB, tf.startGroup)

	tf.AddAction(RB, tf.endGroup)

	tf.AddAction(Add, tf.operator)
	tf.AddAction(Sub, tf.operator)
	tf.AddAction(Mul, tf.operator)
	tf.AddAction(Div, tf.operator)

	chars := "0123456789"
	for _, ch := range chars {
		tf.AddAction(tokenizerEvent(ch), tf.appendDigit)
	}

	return tf
}
