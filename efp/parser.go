package efp

import (
	"fmt"
	"regexp"
	"strconv"
)

type parser struct {
	constructs []construct
	prototype  *ProtoElement
	scope      *Element
	lexer      *lexer
	index      int
	errs       []string
}

func (p *parser) run() {
	for _, c := range p.constructs {
		if c.is(p) {
			c.process(p)
		}
	}
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

func parseDiscoveredAlias(p *parser) {
	alias := p.lexer.tokenString(p.next())
	// go up to find element and add it to the scope
	e := p.prototype
	found := false
	for e != nil && !found {
		if e.fieldAliases != nil && e.fieldAliases[alias] != nil {
			p.addPrototypeField(e.fieldAliases[alias])
			found = true
		}
		if e.elementAliases != nil && e.elementAliases[alias] != nil {
			p.addPrototypeElement(e.elementAliases[alias])
			found = true
		}
		e = e.parent
	}
	if !found {
		p.addError(fmt.Sprintf(errAliasNotVisible, alias, p.prototype.key.key))
	}
}

func parseField(p *parser) {
	f := new(Field)
	f.key = new(Key)
	f.key.key = p.lexer.tokenString(p.next())
	key := p.validateKey(f.key.key)
	p.enforceNext(tknAssign, "Expected '='") // eat =
	f.values = make([]*Value, 0)
	p.parseValue(f.values)
	p.validateField(key, f)
	p.addField(key, f)
}

func (p *parser) parseValue(fv []*Value) {
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parseArrayDeclaration(fv)
		break
	case tknNumber, tknString, tknValue:
		v := new(Value)
		v.value = p.lexer.tokenString(p.next())
		fv = append(fv, v)
		break
	}
}

func (p *parser) parseArrayDeclaration(fv []*Value) {
	current := new(Value)
	p.next() // eat [
	for p.current().tkntype != tknCloseSquare {
		switch p.current().tkntype {
		case tknString, tknValue, tknNumber:
			p.addValueChild(current, p.lexer.tokenString(p.next()))
			break
		case tknOpenSquare:
			p.parseArrayDeclaration(current.values)
		case tknComma:
			p.next()
			break
		default:
			p.addError("Invalid token in array declaration")
			p.next()
			break
		}
	}
	fv = append(fv, current)
	p.next() // eat ]
}

func (p *parser) addTypeChild(t *TypeDeclaration, regex string) {
	if t.types == nil {
		t.types = make([]*TypeDeclaration, 0)
	}
	td := new(TypeDeclaration)
	r, err := regexp.Compile(regex)
	if err == nil {
		p.addError(errInvalidRegex)
	}
	td.value = r
	t.types = append(t.types, td)
}

func (p *parser) addValueChild(fv *Value, regex string) {
	if fv.values == nil {
		fv.values = make([]*Value, 0)
	}
	val := new(Value)
	val.value = regex
	fv.values = append(fv.values, val)
}

func (p *parser) evaluateAlias(alias string) *regexp.Regexp {
	ta := new(TextAlias)
	current := p.prototype
	for current != nil {
		for t, x := range current.textAliases {
			if t == alias {
				ta = &x
			}
		}
		current = current.parent
	}
	if ta == nil {
		p.addError(errAliasNotVisible)
		return nil
	}
	if ta.isRecursive {
		return p.evaluateAlias(ta.value)
	}
	r, err := regexp.Compile(ta.value)
	if err != nil {
		p.addError(errInvalidRegex)
		return nil
	}
	return r
}

func (p *parser) parseTypeDeclaration(t []*TypeDeclaration) {
	switch p.current().tkntype {
	case tknOpenSquare:
		p.parsePrototypeArrayDeclaration(t)
		break
	case tknString:
		td := new(TypeDeclaration)
		r, err := regexp.Compile(p.lexer.tokenString(p.next()))
		if err != nil {
			p.addError(errInvalidRegex)
			return
		}
		td.value = r
		t = append(t, td)
		break
	case tknValue:
		alias := p.lexer.tokenString(p.next())
		td := new(TypeDeclaration)
		td.value = p.evaluateAlias(alias)
	}
	if p.index >= len(p.lexer.tokens) {
		return
	}
	if p.current().tkntype == tknOr {
		p.next()
		p.parseTypeDeclaration(t)
	}
}

func (p *parser) parsePrototypeArrayDeclaration(t []*TypeDeclaration) {
	p.enforceNext(tknOpenSquare, "Expected '['") // eat [
	current := new(TypeDeclaration)
	t = append(t, current)
	current.isArray = true
	if p.current().tkntype == tknNumber {
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		current.min = num
		p.enforceNext(tknColon, "Array minimum must be followed by ':'") // eat ":"
	}
	current.types = make([]*TypeDeclaration, 0)
	p.parseTypeDeclaration(current.types)
	if p.current().tkntype == tknColon {
		p.enforceNext(tknColon, "Array maximum must be preceded by ':'") // eat ":"
		num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
		current.max = num
	}
	p.enforceNext(tknCloseSquare, "Expected ']'") // eat ]
}

func parseElement(p *parser) {
	e := new(Element)
	e.key = new(Key)
	p.parseKey(e.key)
	p.parseParameters(e)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addElement(e.key.key, e)
}

func parseFieldAlias(p *parser) {
	f := new(ProtoField)
	p.next() // eat "alias"
	alias := p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	f.key = new(Key)
	p.parseKey(f.key)
	p.enforceNext(tknColon, "Expected ':'") // eat ":"
	f.types = make([]*TypeDeclaration, 0)
	p.parseTypeDeclaration(f.types)
	p.addFieldAlias(alias, f)
}

func parseElementAlias(p *parser) {
	e := new(ProtoElement)
	p.next() // eat "alias"
	alias := p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	e.key = new(Key)
	p.parseKey(e.key)
	e.parameters = make([]*TypeDeclaration, 0)
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

func (p *parser) parseKeyRegex(k *Key) {
	r, err := regexp.Compile(k.key)
	if err != nil {
		p.addError(fmt.Sprintf(errInvalidRegex, k.key, p.getPrototypeKey()))
	}
	k.regex = r
}

func (p *parser) parseKeyMinimum(k *Key) {
	k.min, _ = strconv.Atoi(p.lexer.tokenString(p.next()))
	p.next() // eat :
}

func (p *parser) parseKeyMaximum(k *Key) {
	k.max, _ = strconv.Atoi(p.lexer.tokenString(p.next()))
	p.next() // eat :
}

func createPrototypeParser(bytes []byte) *parser {
	p := new(parser)
	p.index = 0
	p.importPrototypeConstructs()
	p.lexer = lex(bytes)
	p.prototype = new(ProtoElement)
	p.prototype.key = new(Key)
	p.prototype.key.key = "parent"
	p.prototype.addStandardAliases()
	return p
}

func createPrototypeParserString(data string) *parser {
	return createPrototypeParser([]byte(data))
}

func (p *parser) parseKey(k *Key) {
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
	f := new(ProtoField)
	f.key = new(Key)
	p.parseKey(f.key)
	if f.key.max == 0 {
		f.key.max = 1
	}
	p.enforceNext(tknColon, "Expected ':'") // eat :
	f.types = make([]*TypeDeclaration, 0)
	p.parseTypeDeclaration(f.types)
	p.addPrototypeField(f)
}

func parsePrototypeElement(p *parser) {
	e := new(ProtoElement)
	e.key = new(Key)
	p.parseKey(e.key)
	if e.key.max == 0 {
		e.key.max = 1
	}
	e.parameters = make([]*TypeDeclaration, 0)
	p.parsePrototypeParameters(e.parameters)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addPrototypeElement(e)
}

func (p *parser) parsePrototypeParameters(t []*TypeDeclaration) {
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

func (p *parser) addPrototypeElement(e *ProtoElement) {
	if p.prototype.elements == nil {
		p.prototype.elements = make(map[string]*ProtoElement)
	}
	p.prototype.elements[e.key.key] = e
}

func (p *parser) addPrototypeField(f *ProtoField) {
	if p.prototype.fields == nil {
		p.prototype.fields = make(map[string]*ProtoField)
	}
	p.prototype.fields[f.key.key] = f
}

func (p *parser) addFieldAlias(alias string, f *ProtoField) {
	if p.prototype.fieldAliases == nil {
		p.prototype.fieldAliases = make(map[string]*ProtoField)
	}
	if p.prototype.fieldAliases[alias] != nil {
		p.addError(fmt.Sprintf(errDuplicateAlias, alias, p.prototype.key.key))
	} else {
		p.prototype.fieldAliases[alias] = f
	}
}

func (p *parser) addElementAlias(alias string, e *ProtoElement) {
	if p.prototype.elementAliases == nil {
		p.prototype.elementAliases = make(map[string]*ProtoElement)
	}
	if p.prototype.elementAliases[e.alias] != nil {
		p.addError(fmt.Sprintf(errDuplicateAlias, e.alias, p.prototype.key.key))
	} else {
		p.prototype.elementAliases[e.alias] = e
	}
}

func (p *parser) addElement(key string, e *Element) {
	if p.scope.elements == nil {
		p.scope.elements = make(map[string][]*Element)
	}
	if p.scope.elements[key] == nil {
		p.scope.elements[key] = make([]*Element, 0)
	}
	p.scope.elements[key] = append(p.scope.elements[key], e)
}

func (p *parser) addField(key string, f *Field) {
	if p.scope.fields == nil {
		p.scope.fields = make(map[string][]*Field)
	}
	if p.scope.fields[f.key.key] == nil {
		p.scope.fields[f.key.key] = make([]*Field, 0)
	}
	p.scope.fields[f.key.key] = append(p.scope.fields[f.key.key], f)
}

func (p *parser) validateCompleteElement() {
	for k, v := range p.prototype.fields {
		if v.key.min > len(p.scope.fields[k]) {
			p.addError(fmt.Sprintf(errInsufficientFields, v.key.min, k, p.scope.key.key, len(p.scope.fields[k])))
		} else if v.key.max < len(p.scope.fields[k]) {
			p.addError(fmt.Sprintf(errDuplicateField, v.key.max, k, p.scope.key.key, len(p.scope.fields[k])))
		}
	}
}

func parseElementClosure(p *parser) {
	p.validateCompleteElement()
	p.prototype = p.prototype.parent
	p.scope = p.scope.parent
	p.next()
}

func (p *parser) importValidateConstructs() {
	p.constructs = []construct{
		construct{"field", isField, parseField},
		construct{"element", isElement, parseElement},
		construct{"element closure", isElementClosure, parseElementClosure},
	}
}

func parsePrototypeFieldAlias(p *parser) {
	f := new(ProtoField)
	p.enforceNext(tknValue, "Expected alias keyword") // eat the alias keyword (kw not verified)
	alias := p.lexer.tokenString(p.next())
	p.enforceNext(tknAssign, "Expected '='")
	f.key = new(Key)
	p.parseKey(f.key)
	p.enforceNext(tknAssign, "Expected ':'") // eat =
	f.types = make([]*TypeDeclaration, 0)
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

func (p *parser) parseParameters(e *Element) {
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

func (p *parser) end() {
	p.validateCompleteElement()
}

func (p *parser) validateType(validType *TypeDeclaration, fv *Value) bool {
	if validType.isArray {
		if fv.values == nil {
			return false
		} else {
			if validType.max < len(fv.values) {
				return false // p.addError(fmt.Sprintf(errArrayMaximum, key, p.scope.key.key, prototype.key.min, len(fv.values)))
			} else if validType.min > len(fv.values) {
				return false //	p.addError(fmt.Sprintf(errArrayMinimum, key, p.scope.key.key, prototype.key.min, len(fv.values)))
			}
			for _, v := range fv.values {
				matched := false
				for _, t := range validType.types {
					if p.validateType(t, v) {
						matched = true
						break
					}
				}
				if !matched {
					//p.addError(fmt.Sprintf(errUnmatchedValue, v))
					return false
				}
			}

		}
	} else {
		if fv.values != nil {
			return false
		}
		matched := false
		for _, t := range validType.types {
			if p.validateType(t, fv) {
				matched = true
				break
			}
		}
		if !matched {
			//p.addError(fmt.Sprintf(errUnmatchedValue, v))
			return false
		}
	}
	return true
}

func (p *parser) validateField(key string, f *Field) bool {

	prototype := p.prototype.fields[key]

	for _, v := range f.values {
		matched := false
		for _, t := range prototype.types {
			if p.validateType(t, v) {
				matched = true
			}
		}
		if !matched {
			p.addError(errUnmatchedFieldValue)
			return false // only call out the first one
			// TODO: more specific errors if possible
		}
	}

	return true
}

func (p *parser) addError(err string) {
	if p.errs == nil {
		p.errs = make([]string, 0)
	}
	p.errs = append(p.errs, err)
}
