package protocol

import (
	"errors"
	"testing"
)

func TestProtocolErrorIs(t *testing.T) {
	protoErr := ProtocolError{
		What: FileNameToLarge,
		How:  errors.New("Man..."),
	}
	_ = protoErr

	if !errors.Is(protoErr, FileNameToLarge) {
		t.Fail()
	}

	if !errors.Is(protoErr.What, FileNameToLarge) {
		t.Fail()
	}

	err := errors.New("Super")
	if errors.Is(err, FileNameToLarge) {
		t.Fail()
	}
}
