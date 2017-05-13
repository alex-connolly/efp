package efp

import (
	"fmt"
	"os"
	"strconv"
)

type parser struct {
	constructs []construct
	prototype  *element
	scope      *element
	lexer      *lexer
	index      int
	errs       []string
}

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
	p.parseBytes(bytes)
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

func (p *parser) parseBytes(bytes []byte) {
	p.importParseConstructs()
	p.lexer = lex(bytes)
}

// A construct is a repeated pattern within an efp file
type construct struct {
	name    string // can be used for debugging
	is      func(*parser) bool
	process func(*parser)
}

func (p *parser) peek(offset int) token {
	return p.lexer.tokens[offset]
}

func parseField(p *parser) {
	fmt.Printf("here lads\n")
	f := new(field)
	f.key = p.lexer.tokenString(p.next())
	p.next() // eat =
	f.value = new(fieldValue)
	p.parseFieldValue(f.value)
	p.addField(f)
}

func (fv *fieldValue) addChild(regex string) {
	if fv.children == nil {
		fv.children = make([]*fieldValue, 0)
	}
	val := new(fieldValue)
	val.parent = fv
	val.value = regex
	fv.children = append(fv.children, val)
}

func (p *parser) parseFieldValue(fv *fieldValue) {
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parseArrayDeclaration(fv)
		break
	case tknValue:
		fv.addChild(p.lexer.tokenString(p.next()))
		break
	}
	if p.index >= len(p.lexer.tokens) {
		return
	}
	if p.current().tkntype == tknOr {
		p.next()
		p.parseFieldValue(fv)
	}
}

func (p *parser) parseArrayDeclaration(fv *fieldValue) {
	p.next() // eat [
	fv.isArray = true
	if p.current().tkntype == tknNumber {
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		fv.min = num
		p.next() // eat ":"
	}
	p.parseFieldValue(fv)
	if p.index >= len(p.lexer.tokens) {
		return
	}
	if p.current().tkntype == tknColon {
		p.next() // eat ":"
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		fv.max = num
	}
	p.next() // eat ]
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
	p.next() // eat ":"
	f.key = p.lexer.tokenString(p.next())
	p.next() // eat "="
	f.value = new(fieldValue)
	p.parseFieldValue(f.value)
	p.addFieldAlias(f)
}

func parseElementAlias(p *parser) {
	e := new(element)
	p.next() // eat "alias"
	e.alias = p.lexer.tokenString(p.next())
	p.next() // eat ":"
	e.key = p.lexer.tokenString(p.next())
	p.parsePrototypeParameters(e)
	p.next() // eat "{"
	p.addElementAlias(e)
}

func parsePrototypeField(p *parser) {
	fmt.Printf("here lads\n")
	f := new(field)
	f.key = p.lexer.tokenString(p.next())
	p.next() // eat :
	f.value = new(fieldValue)
	p.parseFieldValue(f.value)
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
		p.prototype.parameters = make([]*fieldValue, 0)
	}
	fv := new(fieldValue)
	p.parseFieldValue(fv)
	p.prototype.parameters = append(p.prototype.parameters, fv)
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
	fmt.Printf("adding field\n")
	if p.prototype.fields == nil {
		p.prototype.fields = make(map[string][]*field)
	}
	if p.prototype.fields[f.key] == nil {
		p.prototype.fields[f.key] = make([]*field, 0)
	}
	p.prototype.fields[f.key] = append(p.prototype.fields[f.key], f)
}

func (p *parser) addFieldAlias(f *field) {
	if p.prototype.declaredFieldAliases == nil {
		p.prototype.declaredFieldAliases = make(map[string][]*field)
	}
	if p.prototype.declaredFieldAliases[f.alias] == nil {
		p.prototype.declaredFieldAliases[f.alias] = make([]*field, 0)
	}
	p.prototype.declaredFieldAliases[f.alias] = append(p.prototype.declaredFieldAliases[f.alias], f)
}

func (p *parser) addElementAlias(e *element) {
	if p.prototype.declaredElementAliases == nil {
		p.prototype.declaredElementAliases = make(map[string][]*element)
	}
	if p.prototype.declaredElementAliases[e.alias] == nil {
		p.prototype.declaredElementAliases[e.alias] = make([]*element, 0)
	}
	p.prototype.declaredElementAliases[e.alias] = append(p.prototype.declaredElementAliases[e.alias], e)
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
		p.errs = append(p.errs, fmt.Sprintf("Duplicate field %s in prototype (max %d)", f.key, 1))
	}
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

func parsePrototypeFieldAlias(p *parser) {
	f := new(field)
	p.next() // eat the alias keyword
	f.key = p.lexer.tokenString(p.next())
	// eat =
	p.next()
	f.value = new(fieldValue)
	p.parseFieldValue(f.value)
	p.addFieldAlias(f)
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
