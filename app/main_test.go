package main

import (
	"testing"
)

func TestSetupRouter(t *testing.T) {
	if setupRouter("test") == nil {
		t.Errorf("setupRouter(%s) got nil, want not nil", "test")
	}
}
