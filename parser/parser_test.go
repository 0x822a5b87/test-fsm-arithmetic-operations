package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParserFsm_Exec(t *testing.T) {
	parserFsm := NewParserFsm("12 + 34 * 56 - (78 - 9)")
	assert.Equal(t, len(parserFsm.tokens), 12)
}

func TestExec(t *testing.T) {
	testCases := []struct {
		data   string
		expect int64
	}{
		{
			"123",
			123,
		},
		{
			"123 + 456",
			123 + 456,
		},
		{
			"123 - 456",
			123 - 456,
		},
		{
			"123 * 456",
			56088,
		},
		{
			"789 / 263",
			3,
		},
		{
			"10 + (12 + 34)",
			56,
		},
		{
			"(12 * 34)",
			12 * 34,
		},
		{
			"(12 + 34 * 56)",
			12 + 34*56,
		},
	}

	for _, testCase := range testCases {
		parseAndExec(t, testCase.data, testCase.expect)
	}
}

func parseAndExec(t *testing.T, data string, expectedValue int64) {
	t.Helper()
	parserFsm := NewParserFsm(data)
	ast := parserFsm.Parse()
	value := ast.Exec()
	assert.Equal(t, expectedValue, value, fmt.Sprintf("data = %s", data))
}
