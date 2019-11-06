package v1

import (
	"testing"
)

func TestInit(t *testing.T) {
	if initValidate() == nil {
		t.Errorf("initValidate() got nil, want not nil")
	}
}
