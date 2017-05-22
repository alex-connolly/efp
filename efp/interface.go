package efp

import (
	"fmt"
	"os"
)

func PrototypeBytes(bytes []byte) (*protoElement, []string) {
	p := new(parser)
	p.index = 0
	p.importPrototypeConstructs()
	p.lexer = lex(bytes)
	p.prototype = new(protoElement)
	p.prototype.addStandardAliases()
	p.run()
	return p.prototype, p.errs
}

func PrototypeString(prototype string) (*protoElement, []string) {
	return PrototypeBytes([]byte(prototype))
}

func PrototypeFile(filename string) (*protoElement, []string) {
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
func (e *protoElement) ValidateBytes(bytes []byte) (*element, []string) {
	p := new(parser)
	p.index = 0
	p.importValidateConstructs()
	p.lexer = lex(bytes)
	p.prototype = e
	p.run()
	return p.scope, p.errs
}

func (p *protoElement) ValidateString(data string) (*element, []string) {
	return p.ValidateBytes([]byte(data))
}

func (p *protoElement) ValidateFile(filename string) (*element, []string) {
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
