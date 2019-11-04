package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	if _, err := Load(); err != nil {
		t.Errorf("TestLoad() got %v, want nil", err)
	}
}
