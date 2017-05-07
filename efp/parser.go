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
	key := p.lexer.tokenString(p.next())
	f.key = key
	p.next() // eat =
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parseFieldArray(f)
		break
	case tknValue:
		p.parseFieldValue(f)
		break
	}
	p.addField(f)
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

func parseElement(p *parser) {

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

func parsePrototypeField(p *parser) {
	f := new(field)
	f.key = p.lexer.tokenString(p.next())
	p.next() // eat :
	switch p.next().tkntype {
	case tknOpenSquare:
		p.parsePrototypeArray(f)
		break
	default:
		p.parsePrototypeFieldValue(f)
		break
	}
	p.addPrototypeField(f)
}

func (p *parser) parsePrototypeFieldValue(f *field) {
	switch p.current().tkntype {
	case tknValue:
		p.addValue(f)
		break
	case tknOpenSquare:
		p.addArrayValue(f)
		break
	default:
		break
	}
	switch p.current().tkntype {
	case tknOr:
		p.next()
		p.parsePrototypeFieldValue(f)
	}
}

func (p *parser) parsePrototypeArray(f *field) {
	switch p.current().tkntype {
	case tknNumber:
		f.min = atoi(p.lexer.tokenString(p.next()))
		break
	case tknValue:

		switch p.current() {
		case tknNumber:
			f.max = atoi(p.lexer.tokenString(p.next()))
		}
		break
	}
	p.next() // eat final ']'
}

func (p *parser) addPrototypeField(f *field) {
	if p.prototype.fields == nil {
		p.prototype.fields = make(map[string][]*field)
	}
	if p.prototype.fields[f.key] == nil {
		p.prototype.fields[f.key] = make([]*field, 0)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Duplicate field in prototype."))
	}
	p.prototype.fields[f.key] = append(p.prototype.fields[f.key], f)
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
	t := p.lexer.tokens[p.index]
	p.index++
	return t
}

func (p *parser) parseParameters() []string {
	var params []string
	for p.current().tkntype != tknCloseBracket {
		params = append(params, p.lexer.tokenString(p.next()))
	}
	if len(params) == 0 {
		return nil
	}
	return params
}

func parsePrototypeElementAlias(p *parser) {
	e := new(element)
	e.key = p.lexer.tokenString(p.next())
	switch p.next().tkntype {
	case tknOpenBrace:
		break
	case tknOpenBracket:
		e.parameters = p.parseParameters()
		p.next() // eat {
		break
	}
	p.addPrototypeElement(e)
}

func (p *parser) addPrototypeElement(e *element) {
	p.prototype.elements[e.key] = append(p.prototype.elements[e.key], e)
}

func parsePrototypeElement(p *parser) {
	e := new(element)
	e.key = p.lexer.tokenString(p.next())
	switch p.next().tkntype {
	case tknOpenBrace:
		break
	case tknOpenBracket:
		p.parseParameters()
		p.next() // eat {
		break
	}
	p.addPrototypeElement(e)
}

func parsePrototypeElementClosure(p *parser) {
	p.scope = p.scope.parent
	p.prototype = p.prototype.parent
	p.index++
}

func (p *parser) importPrototypeConstructs() {
	p.constructs = []construct{
		construct{"prototype field", isPrototypeFieldAlias, parsePrototypeFieldAlias},
		construct{"prototype element", isPrototypeElementAlias, parsePrototypeElementAlias},
		construct{"prototype field", isPrototypeField, parsePrototypeField},
		construct{"prototype element", isPrototypeElement, parsePrototypeElement},
		construct{"prototype element closure", isPrototypeElementClosure, parsePrototypeElementClosure},
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
