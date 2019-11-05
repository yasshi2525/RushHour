package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("fail to get wd %v", err)
	}
	if _, err := Load(dir); err != nil {
		t.Errorf("TestLoad() got %v, want nil", err)
	}
}
