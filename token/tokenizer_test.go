package token

import (
	fsm "0x822a5b87/test-fsm-arithmetic-operations"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tokenFsm := NewTokenizerFsm("")
	if tokenFsm.state != startGroup {
		t.Errorf("expected state [%d], actual state [%d]", startGroup, tokenFsm.state)
	}
}

func TestTokenizer_SingleNumber(t *testing.T) {
	tokenFsm := NewTokenizerFsm("12345")
	peekAndExec(t, tokenFsm)
	assert.Equal(t, appendDigit, tokenFsm.state)
}

func TestTokenizer_TwoNumberOp(t *testing.T) {
	testTwoNumberOp(t, "123 + 456")
	testTwoNumberOp(t, "456 - 123")
	testTwoNumberOp(t, "456 * 123")
	testTwoNumberOp(t, "123456 / 41152")
}

func TestTokenizer_ThreeNumberOp(t *testing.T) {
	tokenFsm := NewTokenizerFsm("123 + 456 * 789")
	states := []fsm.State{
		appendDigit,
		newOperator,
		appendDigit,
		newOperator,
		appendDigit,
		end,
	}
	assertState(t, tokenFsm, states)
}

func TestTokenizer_Parentheses(t *testing.T) {
	tokenFsm := NewTokenizerFsm("(123 + 456) * 789")
	states := []fsm.State{
		startGroup,
		appendDigit,
		newOperator,
		appendDigit,
		endGroup,
		newOperator,
		appendDigit,
		end,
	}
	assertState(t, tokenFsm, states)
}

func TestTokenizer_MulParentheses(t *testing.T) {
	tokenFsm := NewTokenizerFsm("12 + (34 + 56) * (78 / 90)")
	states := []fsm.State{
		appendDigit,
		newOperator,
		startGroup,
		appendDigit,
		newOperator,
		appendDigit,
		endGroup,
		newOperator,
		startGroup,
		appendDigit,
		newOperator,
		appendDigit,
		endGroup,
		end,
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
		appendDigit,
		newOperator,
		appendDigit,
		end,
	}
	assertState(t, tokenFsm, states)
}

func peekAndExec(t *testing.T, fsm *TokenizerFsm) {
	event := fsm.stream.peekEvent()
	fsm.Exec(event)
}
