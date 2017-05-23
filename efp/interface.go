package efp

import (
	"fmt"
	"os"
)

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

// ValidateBytes against the
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

func (p *ProtoElement) ValidateString(data string) (*Element, []string) {
	return p.ValidateBytes([]byte(data))
}

func (p *ProtoElement) ValidateFile(filename string) (*Element, []string) {
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
	return p.ValidateBytes(bytes)
}
