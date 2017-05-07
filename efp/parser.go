package efp

import (
	"fmt"
	"os"
)

func (p *parser) Parse(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open file: file name not found.\n")
		return
	}
	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("Failed to read from file.\n")
		return
	}
	bytes := make([]byte, fi.Size())
	_, err = f.Read(bytes)
	if err != nil {
		fmt.Printf("Failed to read from file.\n")
		return
	}
	p.importParseConstructs()
	p.lexer = lex(bytes)
}

func (p *parser) Prototype(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open file: file name not found.\n")
		return
	}
	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("Failed to read from file.\n")
		return
	}
	bytes := make([]byte, fi.Size())
	_, err = f.Read(bytes)
	if err != nil {
		fmt.Printf("Failed to read from file.\n")
		return
	}
	p.importPrototypeConstructs()
	p.lexer = lex(bytes)
}

// A construct is a repeated pattern within an efp file
type construct struct {
	name    string // can be used for debugging
	is      func(*parser) bool
	process func(*parser)
}

func (p *parser) peek(index int) token {
	return p.lexer.tokens[index]
}

func isField(p *parser) bool {
	return p.peek(p.index).tkntype == tknValue &&
		p.peek(p.index+1).tkntype == tknAssign
}

func isElement(p *parser) bool {
	return p.peek(p.index).tkntype == tknValue &&
		p.peek(p.index+1).tkntype == tknOpenBrace
}

func isElementClosure(p *parser) bool {
	return p.peek(p.index).tkntype == tknCloseBrace
}

func parseField(p *parser) {
	f := new(field)
	f.key = p.lexer.tokenString(p.next())
	p.next() // eat =
	p.parseFieldValue(f)
	p.addField(f)
}

func parseElement(p *parser) {
	e := new(element)
	e.key = p.lexer.tokenString(p.next())
	p.parseParameters()
	p.next() // eat {
	p.addElement(e)
}

func parseFieldAlias(p *parser) {
	f := new(field)
	p.next() // eat "alias"
	f.alias = p.lexer.tokenString(p.next())
	p.next() // eat "="
	f.key = p.lexer.tokenString(p.next())
	p.next() // eat "="
	p.parseFieldValue(f)
	p.addFieldAlias(f)
}

func parseElementAlias(p *parser) {
	e := new(element)
	p.next() // eat "alias"
	e.alias = p.lexer.tokenString(p.next())
	p.next() // eat "="
	e.key = p.lexer.tokenString(p.next())
	p.parsePrototypeParameters(e)
	p.addElementAlias(e)
}

func parsePrototypeField(p *parser) {
	f := new(field)
	f.key = p.lexer.tokenString(p.next())
	p.parseFieldValue(f)
	p.addPrototypeField(f)
}

func parsePrototypeElement(p *parser) {
	e := new(element)
	e.key = p.lexer.tokenString(p.next())
	p.parsePrototypeParameters(e)
	p.next()
	p.addPrototypeElement(e)
}

func (p *parser) parsePrototypeParameters(e *element) {
	// must use current to stop accidentally double-eating the open brace
	if p.current().tkntype != tknOpenBracket {
		return
	}
	p.next() // eat "("
	for p.current().tkntype != tknCloseBracket {
		switch p.current().tkntype {
		case tknComma:
			p.next()
			break
		case tknValue:
			p.parsePrototypeParameter()
			break
		}
	}
	p.next() // eat ")"
}

func (p *parser) parsePrototypeParameter() {
	if p.prototype.parameters == nil {
		p.prototype.parameters = make([]*fieldValue)
	}
	fv := new(fieldValue)

	p.prototype.parameters = append(p.prototype.parameters)

}

func (p *parser) addPrototypeElement(e *element) {
	if p.prototype.elements == nil {
		p.prototype.elements = make(map[string][]*element)
	}
	if p.prototype.elements[e.key] == nil {
		p.prototype.elements[e.key] = make([]*element, 0)
	}
	p.prototype.elements[e.key] = append(p.prototype.elements[e.key], e)
}

func (p *parser) addPrototypeField(f *field) {
	if p.prototype.fields == nil {
		p.prototype.fields = make(map[string][]*field)
	}
	if p.prototype.fields[f.key] == nil {
		p.prototype.fields[f.key] = make([]*field, 0)
	}
	p.prototype.fields[f.key] = append(p.prototype.fields[f.key], f)
}

func (p *parser) addFieldAlias(f *field) {
	if p.scope.fieldAliases == nil {
		p.scope.fieldAliases = make(map[string][]*field)
	}
	if p.scope.fieldAliases[f.key] == nil {
		p.scope.fieldAliases[f.key] = make([]*field, 0)
	}
	p.scope.fieldAliases[f.key] = append(p.scope.fieldAliases[f.key], f)
}

func (p *parser) addElementAlias(e *element) {
	if p.scope.elementAliases == nil {
		p.scope.elementAliases = make(map[string][]*element)
	}
	if p.scope.elementAliases[e.key] == nil {
		p.scope.elementAliases[e.key] = make([]*element, 0)
	}
	p.scope.elementAliases[e.key] = append(p.scope.elementAliases[e.key], e)
}

func (p *parser) addElement(e *element) {
	if p.scope.elements == nil {
		p.scope.elements = make(map[string][]*element)
	}
	if p.scope.elements[e.key] == nil {
		p.scope.elements[e.key] = make([]*element, 0)
	}
	p.scope.elements[e.key] = append(p.scope.elements[e.key], e)
}

func (p *parser) addField(f *field) {
	if p.scope.fields[f.key] == nil {
		p.scope.fields[f.key] = append(p.scope.fields[f.key], f)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Duplicate field %s in prototype (max %d)", f.key, 1))
	}
}

func (p *parser) parseFieldValue(f *field) {
	// first token guaranteed to be a VALUE
	if f.value == nil {
		f.value = make([]*fieldValue, 0)
	}
	fv := new(fieldValue)
	fv.regex = p.lexer.tokenString(p.next())
	f.value = append(f.value, fv)
	switch p.current().tkntype {
	case tknOr:
		p.parseFieldValue(f)
		break
	}
}

func (p *parser) parseFieldArray(f *field) {
	// parse field array
	for p.current().tkntype != tknCloseSquare {
		switch p.next().tkntype {
		case tknComma:
			p.next()
			break
		}
		if f.value == nil {
			f.value = make([]string, 0)
		}
		f.value = append(f.value, p.lexer.tokenString(p.next()))
	}
	p.next() // eat final ']'
}

func parseElementClosure(p *parser) {
	p.prototype = p.prototype.parent
	p.scope = p.scope.parent
	p.next()
}

func (p *parser) importParseConstructs() {
	p.constructs = []construct{
		construct{"field", isField, parseField},
		construct{"element", isElement, parseElement},
		construct{"element closure", isElementClosure, parseElementClosure},
	}
}

func isPrototypeFieldAlias(p *parser) bool {
	return p.lexer.tokenString(p.peek(p.index)) == "alias" &&
		p.peek(p.index+2).tkntype == tknAssign
}

func isPrototypeElementAlias(p *parser) bool {
	return p.lexer.tokenString(p.peek(p.index)) == "alias" &&
		p.peek(p.index+2).tkntype == tknOpenBrace
}

func isPrototypeField(p *parser) bool {
	if len(p.lexer.tokens)-p.index < 2 {
		return false
	}
	return p.peek(p.index+1).tkntype == tknColon
}

func isPrototypeElement(p *parser) bool {
	return p.peek(p.index+1).tkntype == tknOpenBrace
}

func isPrototypeElementClosure(p *parser) bool {
	return p.peek(p.index).tkntype == tknCloseBrace
}

func parsePrototypeFieldAlias(p *parser) {
	f := new(field)
	p.next() // eat the alias keyword
	f.key = p.lexer.tokenString(p.next())
	// eat =
	p.next()
	switch p.next().tkntype {
	case tknOpenSquare:
		// parse field array
		p.parseFieldArray()
		break
	case tknValue:
		p.parseFieldValue()
		break
	}
	p.addPrototypeFieldAlias(f)
}

func (p *parser) addPrototypeFieldAlias(f *field) {
	if p.prototype.fieldAliases[f.key] == nil {
		p.prototype.fieldAliases[f.key] = append(p.prototype.fieldAliases[f.key], f)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Duplicate field in prototype."))
	}
}

func (p *parser) current() token {
	return p.lexer.tokens[p.index]
}

func (p *parser) next() token {
	t := p.current()
	p.index++
	return t
}

func (p *parser) parseParameters() {
	// handle case where no parameters
	if p.current().tkntype != tknOpenBrace {
		return
	}
	p.next() // eat "("
	for p.current().tkntype != tknCloseBracket {
		if p.scope.parameters == nil {
			p.scope.parameters = make([]string, 0)
		}
		p.scope.parameters = append(p.scope.parameters, p.lexer.tokenString(p.next()))
	}
	p.next() // eat ")"'
}

func (p *parser) importPrototypeConstructs() {
	p.constructs = []construct{
		construct{"prototype field", isFieldAlias, parseFieldAlias},
		construct{"prototype element", isElementAlias, parseElementAlias},
		construct{"prototype field", isPrototypeField, parsePrototypeField},
		construct{"prototype element", isPrototypeElement, parsePrototypeElement},
		construct{"element closure", isElementClosure, parseElementClosure},
	}
}

type parser struct {
	constructs []construct
	prototype  *element
	scope      *element
	lexer      *lexer
	tokens     []token
	index      int
	errors     []string
}
