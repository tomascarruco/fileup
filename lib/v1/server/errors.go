package server

import (
	"errors"
	"fmt"
	"reflect"
)

var ()

type ServerError struct {
	What error
	How  string
}

func (se ServerError) Error() string {
	variantName := reflect.TypeOf(se.What).Name()

	return fmt.Sprintf("%s: %s", variantName, se.How)
}

func (se ServerError) Is(target error) bool {
	return errors.Is(target, se.What)
}
