package v1

import "testing"

func TestCustomX(t *testing.T) {
	if CUSTOMX_VALUE != 1 {
		t.Fail()
	}
}
