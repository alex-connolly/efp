package efp

import (
	"fmt"
	"os"
	"regexp"
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
	p.runParserBytes(bytes)
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
	p.runPrototypeBytes(bytes)
}

func (p *parser) run() {
	for _, c := range p.constructs {
		if c.is(p) {
			c.process(p)
		}
	}
}

func (p *parser) runPrototypeBytes(bytes []byte) {
	p.createPrototypeBytes(bytes)
	p.run()
}

func (p *parser) runParserBytes(bytes []byte) {
	p.createPrototypeBytes(bytes)
	p.run()
}

func (p *parser) createPrototypeBytes(bytes []byte) {
	p.importPrototypeConstructs()
	p.lexer = lex(bytes)
	p.prototype = new(element)
}

func (p *parser) createParseBytes(bytes []byte) {
	p.importParseConstructs()
	p.lexer = lex(bytes)
	p.scope = new(element)
}

func (p *parser) createParseString(text string) {
	p.index = 0
	p.createParseBytes([]byte(text))
}

func (p *parser) createPrototypeString(text string) {
	p.index = 0
	p.createPrototypeBytes([]byte(text))
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

func isValue(t tokenType) bool {
	return (t == tknValue) || (t == tknNumber) || (t == tknString)
}

func (p *parser) validateKey(key string) string {
	for k, v := range p.prototype.fields {
		if k == key {
			return k
		} else if v[0].regex != nil {
			if v[0].regex.MatchString(key) {
				return k
			}
		}
	}
	p.addError(fmt.Sprintf("Key %s not matched in prototype element %s", key, p.prototype.key))
	return ""
}

func parseField(p *parser) {
	f := new(field)
	f.key = p.lexer.tokenString(p.next())
	key := p.validateKey(f.key)
	p.enforceNext(tknAssign, "Expected '='") // eat =
	f.value = new(fieldValue)
	p.parseFieldValue(f.value)
	p.validateField(f.value)
	p.addField(key, f)
}

func (p *parser) parseFieldValue(fv *fieldValue) {
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parseArrayDeclaration(fv)
		break
	case tknNumber, tknString, tknValue:
		fv.addChild(p.lexer.tokenString(p.next()))
		break
	}
}

func (p *parser) parseArrayDeclaration(fv *fieldValue) {
	p.next() // eat [
	fv.isArray = true
	for p.current().tkntype != tknCloseSquare {
		switch p.current().tkntype {
		case tknString, tknValue, tknNumber:
			fv.addChild(p.lexer.tokenString(p.next()))
			break
		case tknOpenBracket:
			if fv.children == nil {
				fv.children = make([]*fieldValue, 1)
			}
			fv.children[0] = new(fieldValue)
			p.parseArrayDeclaration(fv.children[0])
		case tknComma:
			p.next()
			break
		default:
			p.addError("Invalid token in array declaration")
			p.next()
			break
		}
	}
	p.next() // eat ]
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

func (p *parser) parsePrototypeFieldValue(fv *fieldValue) {
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parsePrototypeArrayDeclaration(fv)
		break
	case tknValue, tknNumber, tknString:
		fv.addChild(p.lexer.tokenString(p.next()))
		break
	}
	if p.index >= len(p.lexer.tokens) {
		return
	}
	if p.current().tkntype == tknOr {
		p.next()
		p.parsePrototypeFieldValue(fv)
	}
}

func (p *parser) parsePrototypeArrayDeclaration(fv *fieldValue) {
	p.enforceNext(tknOpenSquare, "Expected '['") // eat [
	fv.isArray = true
	if p.current().tkntype == tknNumber {
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		fv.min = num
		p.enforceNext(tknColon, "Array minimum must be followed by ':'") // eat ":"
	}
	p.parsePrototypeFieldValue(fv)
	if p.current().tkntype == tknColon {

		p.enforceNext(tknColon, "Array maximum must be preceded by ':'") // eat ":"
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		fv.max = num
	}
	p.enforceNext(tknCloseSquare, "Expected ']'") // eat ]
}

func parseElement(p *parser) {
	e := new(element)
	e.key = p.lexer.tokenString(p.next())
	p.parseParameters(e)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addElement(e.key, e)
}

func parseFieldAlias(p *parser) {
	f := new(field)
	p.next() // eat "alias"
	f.alias = p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	f.key = p.lexer.tokenString(p.next())
	p.enforceNext(tknColon, "Expected ':'") // eat ":"
	f.value = new(fieldValue)
	p.parsePrototypeFieldValue(f.value)
	p.addFieldAlias(f)
}

func parseElementAlias(p *parser) {
	e := new(element)
	p.next() // eat "alias"
	e.alias = p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	e.key = p.lexer.tokenString(p.next())
	p.parsePrototypeParameters(e)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat "{"
	p.addElementAlias(e)
}

func parsePrototypeField(p *parser) {
	f := new(field)
	f.key = p.lexer.tokenString(p.current())
	if p.next().tkntype == tknString {
		r, err := regexp.Compile(f.key)
		if err != nil {
			p.addError(fmt.Sprintf(errInvalidRegex, f.key, p.prototype.key))
		}
		f.regex = r
	}
	p.enforceNext(tknColon, "Expected ':'") // eat :
	f.value = new(fieldValue)
	p.parsePrototypeFieldValue(f.value)
	p.addPrototypeField(f)
}

func parsePrototypeElement(p *parser) {
	e := new(element)
	e.key = p.lexer.tokenString(p.next())
	p.parsePrototypeParameters(e)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addPrototypeElement(e)
}

func (p *parser) parsePrototypeParameters(e *element) {
	// must use current to stop accidentally double-eating the open brace
	if p.current().tkntype != tknOpenBracket {
		return
	}
	p.enforceNext(tknOpenBracket, "Parameters must open with '('") // eat "("
	for p.current().tkntype != tknCloseBracket {
		switch p.current().tkntype {
		case tknComma:
			p.next()
			break
		case tknValue:
		case tknString:
			p.parsePrototypeParameter()
			break
		}
	}
	p.enforceNext(tknCloseBracket, "Parameters must close with '('") // eat ")"
}

func (p *parser) parsePrototypeParameter() {
	if p.prototype.parameters == nil {
		p.prototype.parameters = make([]*fieldValue, 0)
	}
	fv := new(fieldValue)
	p.parsePrototypeFieldValue(fv)
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
		p.prototype.declaredFieldAliases = make(map[string]*field)
	}
	if p.prototype.declaredFieldAliases[f.alias] != nil {
		p.addError(fmt.Sprintf(errDuplicateAlias, f.alias, p.prototype.key))
	} else {
		p.prototype.declaredFieldAliases[f.alias] = f
	}
}

func (p *parser) addElementAlias(e *element) {
	if p.prototype.declaredElementAliases == nil {
		p.prototype.declaredElementAliases = make(map[string]*element)
	}
	if p.prototype.declaredElementAliases[e.alias] != nil {
		p.addError(fmt.Sprintf(errDuplicateAlias, e.alias, p.prototype.key))
	} else {
		p.prototype.declaredElementAliases[e.alias] = e
	}
}

func (p *parser) addElement(key string, e *element) {
	if p.scope.elements == nil {
		p.scope.elements = make(map[string][]*element)
	}
	if p.scope.elements[key] == nil {
		p.scope.elements[key] = make([]*element, 0)
	}
	p.scope.elements[key] = append(p.scope.elements[key], e)
}

func (p *parser) addField(key string, f *field) {
	if p.scope.fields == nil {
		p.scope.fields = make(map[string][]*field)
	}
	if p.scope.fields[f.key] == nil {
		p.scope.fields[f.key] = make([]*field, 0)
	}
	p.scope.fields[f.key] = append(p.scope.fields[f.key], f)
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
	p.enforceNext(tknValue, "Expected alias keyword") // eat the alias keyword (kw not verified)
	f.key = p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='") // eat =
	f.value = new(fieldValue)
	p.parsePrototypeFieldValue(f.value)
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

func (p *parser) enforceNext(tokType tokenType, err string) token {
	t := p.current()
	p.index++
	if t.tkntype != tokType {
		p.addError(err)
	}
	return t
}

func (p *parser) parseParameters(e *element) {
	// handle case where no parameters
	if p.current().tkntype != tknOpenBrace {
		return
	}
	p.next() // eat "("
	for p.current().tkntype != tknCloseBrace {
		if p.current().tkntype == tknValue {

		} else {

		}
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

func (p *parser) validateField(c *fieldValue) bool {
	return true
}

func (p *parser) addError(err string) {
	if p.errs == nil {
		p.errs = make([]string, 0)
	}
	p.errs = append(p.errs, err)
}
