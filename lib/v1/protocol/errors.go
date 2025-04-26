package protocol

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	FileNameToLargeError     = errors.New("file name to long for header")
	FileNameToSmallError     = errors.New("file name to long for header")
	ChunkLargerThanFileError = errors.New("chunk size larger than total file size")
)

type ProtocolError struct {
	What error
	How  string
}

func (pe ProtocolError) Error() string {
	variantName := reflect.TypeOf(pe.What).Name()

	return fmt.Sprintf("%s: %s", variantName, pe.How)
}

func (pe ProtocolError) Is(target error) bool {
	return errors.Is(target, pe.What)
}
