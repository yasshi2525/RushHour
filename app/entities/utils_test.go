package entities

import (
	"fmt"
	"testing"
)

// TestCase represents Op(input) should be want.
type TestCase struct {
	Name  string
	Input interface{}
	Want  interface{}
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
		t.Run(c.Name, func(t *testing.T) {
			if got := op(c.Input); got != c.Want {
				t.Errorf("want %v, but %v got %v", c.Want, c.Input, got)
			}
		})
	}
}

// TestCaseLineTask represents assertion of LineTask type and base(moving/staying) id
type TestCaseLineTask struct {
	Name string
	Type LineTaskType
	Base interface{}
}

// TestCaseLineTasks is set of test case
type TestCaseLineTasks []TestCaseLineTask

// Assert walk through loop of LineTask and assert condition.
func (cases TestCaseLineTasks) Assert(t *testing.T, lt *LineTask) {
	for i, c := range cases {
		t.Run(fmt.Sprintf("[%d] %s", i, c.Name), func(t *testing.T) {
			if lt == (*LineTask)(nil) {
				t.Error("lt want not be nil, but got nil")
			}
			if lt.TaskType != c.Type {
				t.Errorf("type error: want %v, but got %v", c.Type, lt.TaskType)
			}
			var base interface{}
			switch c.Type {
			case OnDeparture:
				base = lt.Stay
			case OnMoving:
				fallthrough
			case OnStopping:
				fallthrough
			case OnPassing:
				base = lt.Moving
			}
			if base != c.Base {
				t.Errorf("id error: %v want %v, but got %v", c.Type, c.Base, base)
			}
			lt = lt.next
		})
	}
}
