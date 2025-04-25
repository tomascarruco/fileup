package protocol

import (
	"errors"
	"fmt"
	"reflect"
)

type ProtocolErrorVariant uint8

const (
	FileNameToLarge ProtocolErrorVariant = iota
)

func (pev ProtocolErrorVariant) Error() string {
	return reflect.TypeOf(pev).Name()
}

func (pev ProtocolErrorVariant) Is(err error) bool {
	if reflect.TypeOf(err) == reflect.TypeOf(ProtocolError{}) {
		return err.(ProtocolError).What == pev
	}
	return false
}

type ProtocolError struct {
	What ProtocolErrorVariant
	How  error
}

func (pe ProtocolError) Error() string {
	variantName := reflect.TypeOf(pe.What).Name()

	return fmt.Sprintf("%s: %s", variantName, pe.How.Error())
}

func (pe ProtocolError) Is(target error) bool {
	return errors.Is(target, pe.What)
}
