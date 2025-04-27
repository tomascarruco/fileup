package protocol

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	FileNameToLargeError = errors.New("file name to long for header")
	FileNameToSmallError = errors.New("file name to small for header")

	FolderNameToLargeError = errors.New("folder name to long for header")
	FolderPathToLargeError = errors.New("folder path to long for header")
	FolderNameToSmallError = errors.New("folder name to small for header")
	FolderPathToSmallError = errors.New("folder path to small for header")

	InvalidLocalFolderAsPathError = errors.New("Cannot use '.' as a path")

	ChunkLargerThanFileError = errors.New("chunk size larger than total file size")
	ReadLargerThenChunkError = errors.New("content cannot feat into chunk")
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
