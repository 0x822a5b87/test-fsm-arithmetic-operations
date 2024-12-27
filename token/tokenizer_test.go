package token

import (
	fsm "0x822a5b87/test-fsm-arithmetic-operations"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tokenFsm := NewTokenizerFsm("")
	if tokenFsm.state != StartGroup {
		t.Errorf("expected state [%d], actual state [%d]", StartGroup, tokenFsm.state)
	}
}

func TestTokenizer_SingleNumber(t *testing.T) {
	tokenFsm := NewTokenizerFsm("12345")
	peekAndExec(t, tokenFsm)
	assert.Equal(t, AppendDigit, tokenFsm.state)
}

func TestTokenizer_TwoNumberOp(t *testing.T) {
	testTwoNumberOp(t, "123 + 456")
	testTwoNumberOp(t, "456 - 123")
	testTwoNumberOp(t, "456 * 123")
	testTwoNumberOp(t, "123456 / 41152")
}

func TestTokenizer_Parentheses(t *testing.T) {
	tokenFsm := NewTokenizerFsm("(123 + 456) * 789")
	states := []fsm.State{
		StartGroup,
		AppendDigit,
		NewOperator,
		AppendDigit,
		EndGroup,
		NewOperator,
		AppendDigit,
		End,
	}
	assertState(t, tokenFsm, states)
}

func TestTokenizer_MulParentheses(t *testing.T) {
	tokenFsm := NewTokenizerFsm("12 + (34 + 56) * (78 / 90)")
	states := []fsm.State{
		AppendDigit,
		NewOperator,
		StartGroup,
		AppendDigit,
		NewOperator,
		AppendDigit,
		EndGroup,
		NewOperator,
		StartGroup,
		AppendDigit,
		NewOperator,
		AppendDigit,
		EndGroup,
		End,
	}
	assertState(t, tokenFsm, states)
}

func assertState(t *testing.T, fsm *TokenizerFsm, states []fsm.State) {
	for _, state := range states {
		peekAndExec(t, fsm)
		assert.Equal(t, state, fsm.state)
	}
}

func testTwoNumberOp(t *testing.T, input string) {
	tokenFsm := NewTokenizerFsm(input)
	states := []fsm.State{
		AppendDigit,
		NewOperator,
		AppendDigit,
		End,
	}
	assertState(t, tokenFsm, states)
}

func peekAndExec(t *testing.T, fsm *TokenizerFsm) {
	event := fsm.stream.peekEvent()
	fsm.Exec(event)
}
