package entities

import "testing"

// TestCase represents Op(input) should be want.
type TestCase struct {
	name  string
	input interface{}
	want  interface{}
}

// TestCases is set of test case
type TestCases []TestCase

// Assert assert all test case satisfy condition.
func (cases TestCases) Assert(t *testing.T, callback ...func(interface{}) interface{}) {
	op := func(input interface{}) interface{} {
		return input
	}
	if len(callback) > 0 {
		op = callback[0]
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := op(c.input); got != c.want {
				t.Errorf("want = %v, but %v got %v", c.want, c.input, got)
			}
		})
	}
}
