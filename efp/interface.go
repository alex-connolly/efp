package efp

import (
	"fmt"
	"os"
)

// PrototypeBytes runs the prototype parser on a byte array
func PrototypeBytes(bytes []byte) (*ProtoElement, []string) {
	p := createPrototypeParser(bytes)
	p.run()
	return p.prototype, p.errs
}

// PrototypeString forms a prototype element from an input string
func PrototypeString(prototype string) (*ProtoElement, []string) {
	return PrototypeBytes([]byte(prototype))
}

// PrototypeFile forms a prototype element from an input file
func PrototypeFile(filename string) (*ProtoElement, []string) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to open file: file name %s not found.", filename)}
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to read from file %s.", filename)}
	}
	bytes := make([]byte, fi.Size())
	_, err = f.Read(bytes)
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to read from file %s.", filename)}
	}
	return PrototypeBytes(bytes)
}

// ValidateBytes against prototype
func (e *ProtoElement) ValidateBytes(bytes []byte) (*Element, []string) {
	p := new(parser)
	p.index = 0
	p.importValidateConstructs()
	p.lexer = lex(bytes)
	p.prototype = e

	p.scope = new(Element)
	p.scope.key = new(Key)
	p.scope.key.key = "parent"
	p.run()
	p.end()
	return p.scope, p.errs
}

// ValidateString against prototype
func (e *ProtoElement) ValidateString(data string) (*Element, []string) {
	return e.ValidateBytes([]byte(data))
}

// ValidateFile validates a file by filename
func (e *ProtoElement) ValidateFile(filename string) (*Element, []string) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to open file: file name %s not found.", filename)}
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to read from file %s.", filename)}
	}
	bytes := make([]byte, fi.Size())
	_, err = f.Read(bytes)
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to read from file %s.", filename)}
	}
	return e.ValidateBytes(bytes)
}

// ValidateFiles validates a series of files with the same prototype element
func (e *ProtoElement) ValidateFiles(filenames []string) ([]*Element, [][]string) {
	elements := make([]*Element, len(filenames))
	errors := make([][]string, len(filenames))
	for i, name := range filenames {
		elements[i], errors[i] = e.ValidateFile(name)
	}
	return elements, errors
}
