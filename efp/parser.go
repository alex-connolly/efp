package efp

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type parser struct {
	constructs []construct
	prototype  *protoElement
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
	p.prototype = new(protoElement)
	p.prototype.addStandardAliases()
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
		} else if v != nil && v.key != nil && v.key.regex != nil {
			if v.key.regex.MatchString(key) {
				return k
			}
		}
	}
	if p.prototype.key == nil {
		p.addError(fmt.Sprintf("Key %s not matched in global scope", key))
	} else {
		p.addError(fmt.Sprintf("Key %s not matched in prototype element %s", key, p.prototype.key.key))
	}
	return ""
}

func parseField(p *parser) {
	f := new(field)
	f.key = new(key)
	f.key.key = p.lexer.tokenString(p.next())
	key := p.validateKey(f.key.key)
	p.enforceNext(tknAssign, "Expected '='") // eat =
	f.value = new(value)
	p.parsevalue(f.value)
	p.validateField(key, f)
	p.addField(key, f)
}

func (p *parser) parsevalue(fv *value) {
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parseArrayDeclaration(fv)
		break
	case tknNumber, tknString, tknValue:
		fv.addChild(p.lexer.tokenString(p.next()))
		break
	}
}

func (p *parser) parseArrayDeclaration(fv *value) {
	p.next() // eat [
	for p.current().tkntype != tknCloseSquare {
		switch p.current().tkntype {
		case tknString, tknValue, tknNumber:
			fv.addChild(p.lexer.tokenString(p.next()))
			break
		case tknOpenBracket:
			if fv.children == nil {
				fv.children = make([]*value, 1)
			}
			fv.children[0] = new(value)
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

func (t *typeDeclaration) addChild(regex string) {
	if t.types == nil {
		t.types = make([]*typeDeclaration, 0)
	}
	td := new(typeDeclaration)
	td.value = regex
	t.types = append(t.types, td)
}

func (fv *value) addChild(regex string) {
	if fv.children == nil {
		fv.children = make([]*value, 0)
	}
	val := new(value)
	val.parent = fv
	val.value = regex
	fv.children = append(fv.children, val)
}

func (p *parser) parseTypeDeclaration(t []*typeDeclaration) {
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parsePrototypeArrayDeclaration(t)
		break
	case tknValue, tknNumber, tknString:
		td := new(typeDeclaration)
		td.value = p.lexer.tokenString(p.next())
		t = append(t, td)
		break
	}
	if p.index >= len(p.lexer.tokens) {
		return
	}
	if p.current().tkntype == tknOr {
		p.next()
		p.parseTypeDeclaration(t)
	}
}

func (p *parser) parsePrototypeArrayDeclaration(t []*typeDeclaration) {
	p.enforceNext(tknOpenSquare, "Expected '['") // eat [
	if p.current().tkntype == tknNumber {
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		t[len(t)-1].min = num
		p.enforceNext(tknColon, "Array minimum must be followed by ':'") // eat ":"
	}
	p.parseTypeDeclaration(t)
	if p.current().tkntype == tknColon {
		p.enforceNext(tknColon, "Array maximum must be preceded by ':'") // eat ":"
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		t[len(t)-1].max = num
	}
	p.enforceNext(tknCloseSquare, "Expected ']'") // eat ]
}

func parseElement(p *parser) {
	e := new(element)
	e.key = new(key)
	p.parseKey(e.key)
	p.parseParameters(e)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addElement(e.key.key, e)
}

func parseFieldAlias(p *parser) {
	f := new(protoField)
	p.next() // eat "alias"
	alias := p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	f.key = new(key)
	p.parseKey(f.key)
	p.enforceNext(tknColon, "Expected ':'") // eat ":"
	f.types = make([]*typeDeclaration, 0)
	p.parseTypeDeclaration(f.types)
	p.addFieldAlias(alias, f)
}

func parseElementAlias(p *parser) {
	e := new(protoElement)
	p.next() // eat "alias"
	alias := p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	e.key = new(key)
	p.parseKey(e.key)
	e.parameters = make([]*typeDeclaration, 0)
	p.parsePrototypeParameters(e.parameters)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat "{"
	p.addElementAlias(alias, e)
}

func (p *parser) getPrototypeKey() string {
	if p.prototype.key == nil {
		return "global"
	}
	return p.prototype.key.key
}

func (p *parser) parseKeyRegex(k *key) {
	r, err := regexp.Compile(k.key)
	if err != nil {
		p.addError(fmt.Sprintf(errInvalidRegex, k.key, p.getPrototypeKey()))
	}
	k.regex = r
}

func (p *parser) parseKeyMinimum(k *key) {
	k.min, _ = strconv.Atoi(p.lexer.tokenString(p.next()))
	p.next() // eat :
}

func (p *parser) parseKeyMaximum(k *key) {
	k.max, _ = strconv.Atoi(p.lexer.tokenString(p.next()))
	p.next() // eat :
}

func (p *parser) parseKey(k *key) {
	switch p.current().tkntype {
	case tknValue:
		k.key = p.lexer.tokenString(p.next())
		break
	case tknString:
		k.key = p.lexer.tokenString(p.next())
		p.parseKeyRegex(k)
		break
	case tknOpenCorner:
		p.next() // eat <
		switch p.current().tkntype {
		case tknNumber:
			p.parseKeyMinimum(k)
			switch p.current().tkntype {
			case tknValue:
				k.key = p.lexer.tokenString(p.next())
				switch p.current().tkntype {
				case tknColon:
					p.parseKeyMaximum(k)
					break
				case tknCloseCorner:
					break
				}
				break
			case tknString:
				k.key = p.lexer.tokenString(p.next())
				p.parseKeyRegex(k)
				switch p.current().tkntype {
				case tknColon:
					p.parseKeyMaximum(k)
					break
				case tknCloseCorner:
					break
				}
				break
			}
			break
		case tknValue:
			k.key = p.lexer.tokenString(p.next())
			break
		case tknString:
			k.key = p.lexer.tokenString(p.next())
			p.parseKeyRegex(k)
			break
		}
		p.enforceNext(tknCloseCorner, "Open corner in field key must be closed") // eat >
		break
	}
}

func parsePrototypeField(p *parser) {
	f := new(protoField)
	f.key = new(key)
	p.parseKey(f.key)
	p.enforceNext(tknColon, "Expected ':'") // eat :
	f.types = make([]*typeDeclaration, 0)
	p.parseTypeDeclaration(f.types)
	p.addPrototypeField(f)
}

func parsePrototypeElement(p *parser) {
	e := new(protoElement)
	e.key = new(key)
	p.parseKey(e.key)
	e.parameters = make([]*typeDeclaration, 0)
	p.parsePrototypeParameters(e.parameters)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addPrototypeElement(e)
}

func (p *parser) parsePrototypeParameters(t []*typeDeclaration) {
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
		case tknValue, tknString:
			p.parseTypeDeclaration(t)
			break
		}
	}
	p.enforceNext(tknCloseBracket, "Parameters must close with ')'") // eat ")"
}

func (p *parser) addPrototypeElement(e *protoElement) {
	if p.prototype.elements == nil {
		p.prototype.elements = make(map[string]*protoElement)
	}
	p.prototype.elements[e.key.key] = e
}

func (p *parser) addPrototypeField(f *protoField) {
	if p.prototype.fields == nil {
		p.prototype.fields = make(map[string]*protoField)
	}
	p.prototype.fields[f.key.key] = f
}

func (p *parser) addFieldAlias(alias string, f *protoField) {
	if p.prototype.fieldAliases == nil {
		p.prototype.fieldAliases = make(map[string]*protoField)
	}
	if p.prototype.fieldAliases[alias] != nil {
		p.addError(fmt.Sprintf(errDuplicateAlias, alias, p.prototype.key.key))
	} else {
		p.prototype.fieldAliases[alias] = f
	}
}

func (p *parser) addElementAlias(alias string, e *protoElement) {
	if p.prototype.elementAliases == nil {
		p.prototype.elementAliases = make(map[string]*protoElement)
	}
	if p.prototype.elementAliases[e.alias] != nil {
		p.addError(fmt.Sprintf(errDuplicateAlias, e.alias, p.prototype.key.key))
	} else {
		p.prototype.elementAliases[e.alias] = e
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
	if p.scope.fields[f.key.key] == nil {
		p.scope.fields[f.key.key] = make([]*field, 0)
	}
	p.scope.fields[f.key.key] = append(p.scope.fields[f.key.key], f)
}

func (p *parser) validateCompleteElement() {
	fmt.Printf("hi\n")
	for k, v := range p.prototype.fields {
		if v.key.min > len(p.scope.fields[k]) {
			p.addError(errInsufficientFields)
		}
	}
}

func parseElementClosure(p *parser) {
	fmt.Printf("hip\n")
	p.validateCompleteElement()
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
	f := new(protoField)
	p.enforceNext(tknValue, "Expected alias keyword") // eat the alias keyword (kw not verified)
	alias := p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='")
	f.key = new(key)
	p.parseKey(f.key)
	p.enforceNext(tknAssign, "Expected ':'") // eat =
	f.types = make([]*typeDeclaration, 0)
	p.parseTypeDeclaration(f.types)
	p.addFieldAlias(alias, f)
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

func (p *parser) validateField(key string, f *field) bool {

	// check field frequency
	if p.prototype.fields[key].key.max != 0 {
		if len(p.scope.fields[key]) >= p.prototype.fields[key].key.max {
			p.addError(errDuplicateField)
			return false
		}
	}

	/* check single value regex matching
	if len(f.value.children) == 1 {
		fmt.Println("here")
		if p.prototype.fields[key][0].value.children[0].value != nil {
			fmt.Println("and here")
			if !p.prototype.fields[key][0].key.regex.MatchString(f.value.children[0].value) {
				p.addError(errInvalidRegex)
				return false
			}
		}
	}*/

	return true
}

func (p *parser) addError(err string) {
	if p.errs == nil {
		p.errs = make([]string, 0)
	}
	p.errs = append(p.errs, err)
}
