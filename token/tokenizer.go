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

	Null    TokenizerEvent = 0
	Space   TokenizerEvent = ' '
	Tab     TokenizerEvent = '\t'
	Return  TokenizerEvent = '\r'
	NewLine TokenizerEvent = '\n'
	LB      TokenizerEvent = '('
	RB      TokenizerEvent = ')'
	Add     TokenizerEvent = '+'
	Sub     TokenizerEvent = '-'
	Mul     TokenizerEvent = '*'
	Div     TokenizerEvent = '/'
)

type TokenizerEvent byte

func (t TokenizerEvent) EventType() int {
	return int(t)
}

type TokenizerFsm struct {
	state  fsm.State
	acts   map[TokenizerEvent]fsm.Action
	stream charStream
	tokens []string
}

func (fsm *TokenizerFsm) Exec(event TokenizerEvent) {
	action, ok := fsm.acts[event]
	if !ok {
		panic(fmt.Errorf("action not found for [%d]", event))
	}
	action(event)
}

func (fsm *TokenizerFsm) AddAction(event TokenizerEvent, action fsm.Action) {
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
	event = event.(TokenizerEvent)
	fsm.stream.nextEvent()
	fsm.tokens = append(fsm.tokens, "(")
}

func (fsm *TokenizerFsm) endGroup(event fsm.Event) {
	fsm.setState(EndGroup)
	event = event.(TokenizerEvent)
	fsm.stream.nextEvent()
	fsm.tokens = append(fsm.tokens, ")")
}

func (fsm *TokenizerFsm) end(event fsm.Event) {
	fsm.setState(End)
}

func (fsm *TokenizerFsm) operator(event fsm.Event) {
	fsm.setState(NewOperator)
	fsm.stream.nextEvent()

	e := event.(TokenizerEvent)
	buff := bytes.NewBufferString("")
	buff.WriteByte(byte(e))
	fsm.tokens = append(fsm.tokens, buff.String())
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

	fsm.tokens = append(fsm.tokens, buff.String())
}

func NewTokenizerFsm(data string) *TokenizerFsm {
	tf := &TokenizerFsm{
		state:  StartGroup,
		acts:   make(map[TokenizerEvent]fsm.Action),
		stream: charStream{data: data},
		tokens: make([]string, 0),
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
		tf.AddAction(TokenizerEvent(ch), tf.appendDigit)
	}

	return tf
}
