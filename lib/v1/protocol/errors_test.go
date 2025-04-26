package protocol

import (
	"errors"
	"testing"
)

func TestProtocolErrorIs(t *testing.T) {
	protoErr := ProtocolError{
		What: FileNameToLargeError,
		How:  "Man...",
	}

	if !errors.Is(protoErr, FileNameToLargeError) {
		t.Fail()
	}

	if !errors.Is(protoErr.What, FileNameToLargeError) {
		t.Fail()
	}

	err := errors.New("Super")
	if errors.Is(err, FileNameToLargeError) {
		t.Fail()
	}
}
