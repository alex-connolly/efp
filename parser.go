package efp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const parentKey = "parent"

type parser struct {
	constructs []construct
	prototype  *ProtoElement
	scope      *Element
	lexer      *lexer
	index      int
	errs       []string
}

func (p *parser) run() {
	if p.index >= len(p.lexer.tokens) {
		return
	}
	found := false
	for _, c := range p.constructs {
		if c.is(p) {
			//fmt.Printf("FOUND: %s at index %d\n", c.name, p.index)
			c.process(p)
			found = true
			break
		}
	}
	if !found {
		p.addError(fmt.Sprintf(errUnrecognisedConstruct, p.lexer.tokenString(p.current())))
		p.next()
	}
	p.run()
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
	p.addError(fmt.Sprintf("Key %s not matched in prototype element %s", key, p.prototype.key.key))
	return ""
}

func parseDiscoveredAlias(p *parser) {
	alias := p.lexer.tokenString(p.next())
	//fmt.Printf("alias: %s [%d]\n", alias, len(alias))
	// go up to find element and add it to the scope
	e := p.prototype
	found := false
	for e != nil && !found {
		if e.fieldAliases != nil && e.fieldAliases[alias] != nil {
			p.addPrototypeField(e.fieldAliases[alias])
			found = true
			break
		}
		if e.elementAliases != nil && e.elementAliases[alias] != nil {
			p.addPrototypeElement(e.elementAliases[alias])
			found = true
			break
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
	f.key.key = strval(p.lexer.tokenString(p.next()))
	key := p.validateKey(f.key.key)
	p.enforceNext(tknAssign, "Expected '='") // eat =
	f.values = make([]*Value, 0)
	f.values = p.parseValue(f.values)
	p.validateField(key, f)
	p.addField(key, f)
}

func (p *parser) parseValue(fv []*Value) []*Value {
	switch p.current().tkntype {
	case tknOpenSquare:
		fv = p.parseArrayDeclaration(fv)
		break
	case tknNumber, tknString, tknValue:
		v := new(Value)
		v.value = p.lexer.tokenString(p.next())
		fv = append(fv, v)
		break
	}
	return fv
}

func (p *parser) parseArrayDeclaration(fv []*Value) []*Value {
	current := new(Value)
	p.next() // eat [
	for p.current().tkntype != tknCloseSquare {
		switch p.current().tkntype {
		case tknString, tknValue, tknNumber:
			value := p.lexer.tokenString(p.next())
			p.addValueChild(current, value)
			break
		case tknOpenSquare:
			if current.values == nil {
				current.values = make([]*Value, 0)
			}
			current.values = p.parseArrayDeclaration(current.values)
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
	return fv
}

func (p *parser) addValueChild(fv *Value, data string) {
	if fv.values == nil {
		fv.values = make([]*Value, 0)
	}
	val := new(Value)
	val.value = data
	fv.values = append(fv.values, val)
}

func (p *parser) findTextAlias(alias string) *TextAlias {
	current := p.prototype
	for current != nil {
		for t, x := range current.textAliases {
			if t == alias {
				return &x
			}
		}
		current = current.parent
	}
	return nil
}

func (p *parser) evaluateAlias(alias string) *regexp.Regexp {
	ta := p.findTextAlias(alias)
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

func (p *parser) parseTypeDeclaration(t []*TypeDeclaration) []*TypeDeclaration {
	switch p.current().tkntype {
	case tknOpenSquare:
		t = p.parsePrototypeArrayDeclaration(t)
		break
	case tknString:
		td := new(TypeDeclaration)
		r, err := regexp.Compile(strval(p.lexer.tokenString(p.next())))
		if err != nil {
			p.addError(errInvalidRegex)
			return t
		}
		td.value = r
		t = append(t, td)
		break
	case tknValue:
		alias := strval(p.lexer.tokenString(p.next()))
		td := new(TypeDeclaration)
		td.value = p.evaluateAlias(alias)
		t = append(t, td)
	default:
		//fmt.Printf("wrong token: %d @ %d\n", p.current().tkntype, p.index)
	}
	if p.index >= len(p.lexer.tokens) {
		return t
	}
	if p.current().tkntype == tknOr {
		p.next() // eat |
		t = p.parseTypeDeclaration(t)
	}
	return t
}

func (p *parser) parsePrototypeArrayDeclaration(t []*TypeDeclaration) []*TypeDeclaration {
	p.enforceNext(tknOpenSquare, "Expected '['") // eat [
	current := new(TypeDeclaration)
	current.isArray = true
	t = append(t, current)
	// second condition dodges case:
	// what about [2:string] vs [string:2] (remember [LIMIT:string] and [string:LIMIT])
	// no method of differentiation!
	if p.lexer.tokens[p.index+1].tkntype == tknColon && p.lexer.tokens[p.index-1].tkntype == tknOpenSquare {
		switch p.current().tkntype {
		case tknNumber:
			num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
			current.min = num
			p.enforceNext(tknColon, "Array minimum must be followed by ':'") // eat ":"
		case tknValue:
			a := p.findTextAlias(p.lexer.tokenString(p.current()))
			i, err := strconv.Atoi(a.value)
			if err != nil {
				// this is a max, not a min, so ignore
			} else {
				p.next()
				current.min = i
				p.enforceNext(tknColon, "Array minimum must be followed by ':'") // eat ":"
			}
		default:
			// ignore
			//fmt.Printf("wrong token: %d?\n", p.current().tkntype)
		}

	}
	current.types = make([]*TypeDeclaration, 0)
	current.types = p.parseTypeDeclaration(current.types)
	if p.current().tkntype == tknColon {
		p.enforceNext(tknColon, "Array maximum must be preceded by ':'") // eat ":"
		switch p.current().tkntype {
		case tknNumber:
			num, _ := strconv.Atoi(p.lexer.tokenString(p.next()))
			current.max = num
		case tknValue:
			a := p.findTextAlias(p.lexer.tokenString(p.next()))
			i, err := strconv.Atoi(a.value)
			if err != nil {
				p.addError(errInvalidLimitAlias)
			} else {
				current.max = i
			}
		}

	}
	p.enforceNext(tknCloseSquare, "Expected ']'") // eat ]
	return t
}

func parseElement(p *parser) {
	e := new(Element)
	e.key = new(Key)
	p.parseKey(e.key)
	p.parseParameters(e)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addElement(e.key.key, e)
	e.parent = p.scope
	p.prototype = p.prototype.Element(e.key.key)
	p.scope = e
}

func parseFieldAlias(p *parser) {
	f := new(ProtoField)
	p.next() // eat "alias"
	alias := strval(p.lexer.tokenString(p.next()))
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	f.key = new(Key)
	p.parseKey(f.key)
	p.enforceNext(tknColon, "Expected ':'") // eat ":"
	f.types = make([]*TypeDeclaration, 0)
	f.types = p.parseTypeDeclaration(f.types)
	p.addFieldAlias(alias, f)
}

func parseElementAlias(p *parser) {
	e := new(ProtoElement)
	p.next() // eat "alias"
	alias := strval(p.lexer.tokenString(p.next()))
	p.enforceNext(tknAssign, "Expected '='") // eat "="
	e.key = new(Key)
	p.parseKey(e.key)
	e.parameters = make([]*TypeDeclaration, 0)
	e.parameters = p.parsePrototypeParameters(e.parameters)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat "{"
	p.addElementAlias(alias, e)
	e.parent = p.prototype
	p.prototype = e
}

func (p *parser) getPrototypeKey() string {
	if p.prototype.key == nil {
		return parentKey
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

func createPrototypeParser(bytes []byte) *parser {
	p := new(parser)
	p.index = 0
	p.importPrototypeConstructs()
	p.lexer = lex(bytes)
	// add all errors (possible to improve process in future)
	if p.lexer.errors != nil && len(p.lexer.errors) > 0 {
		// MUST leave this condition in --> 0 errors should be nil, not an empty array
		p.errs = make([]string, 0)
		for _, err := range p.lexer.errors {
			p.errs = append(p.errs, err)
		}
	}
	p.prototype = new(ProtoElement)
	p.prototype.key = new(Key)
	p.prototype.key.key = parentKey
	p.prototype.addStandardAliases()
	return p
}

func (p *parser) isAliasAvailable(alias string) bool {
	if p.prototype.textAliases != nil {
		_, found := p.prototype.textAliases[alias]
		if found {
			return false
		}
	}
	if p.prototype.fieldAliases != nil {
		if p.prototype.fieldAliases[alias] != nil {
			return false
		}
	}
	if p.prototype.elementAliases != nil {
		if p.prototype.elementAliases[alias] != nil {
			return false
		}
	}
	return true
}

func createPrototypeParserString(data string) *parser {
	return createPrototypeParser([]byte(data))
}

func (p *parser) parseKeyMaximum(k *Key) {
	if p.current().tkntype == tknColon {
		p.next() // eat :
		switch p.current().tkntype {
		case tknValue:
			a := p.findTextAlias(p.lexer.tokenString(p.next()))
			i, err := strconv.Atoi(a.value)
			if err != nil {
				p.addError(errInvalidLimitAlias)
			}
			k.max = i
		case tknNumber:
			k.max, _ = strconv.Atoi(p.lexer.tokenString(p.next()))
		}
	}
}

func (p *parser) parseKeyMinimum(k *Key) {
	if p.lexer.tokens[p.index+1].tkntype == tknColon {
		if p.current().tkntype == tknValue {
			a := p.findTextAlias(p.lexer.tokenString(p.next()))
			if a == nil {
				// do someting
				return
			}
			i, err := strconv.Atoi(a.value)
			if err != nil {
				p.addError(errInvalidLimitAlias)
			} else {
				k.min = i
				p.next() // eat :
			}
		} else if p.current().tkntype == tknNumber {
			k.min, _ = strconv.Atoi(p.lexer.tokenString(p.next()))
			p.next() // eat :
		}

	}
}

func (p *parser) parseKeyText(k *Key) {
	switch p.current().tkntype {
	case tknString:
		k.key = strval(p.lexer.tokenString(p.next()))
		p.parseKeyRegex(k)
	case tknValue:
		k.key = p.lexer.tokenString(p.next())
		//default:
		//fmt.Printf("wrong token: %d\n", p.current().tkntype)
	}
	//fmt.Printf("parsed key: %s\n", k.key)
}

func (p *parser) parseKey(k *Key) {
	if p.current().tkntype == tknOpenCorner {
		p.next()             // eat <
		p.parseKeyMinimum(k) // can handle it being optional
		p.parseKeyText(k)
		p.parseKeyMaximum(k) // can handle it being optional
		p.next()             // eat >
	} else {
		p.parseKeyText(k)
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
	f.types = p.parseTypeDeclaration(f.types)
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
	e.parameters = p.parsePrototypeParameters(e.parameters)
	p.enforceNext(tknOpenBrace, "Expected '{'") // eat {
	p.addPrototypeElement(e)
	e.parent = p.prototype
	p.prototype = e
}

func (p *parser) parsePrototypeParameters(t []*TypeDeclaration) []*TypeDeclaration {
	if p.current().tkntype != tknOpenBracket {
		return t
	}
	p.enforceNext(tknOpenBracket, "Parameters must open with '('") // eat "("
	for p.current().tkntype != tknCloseBracket {
		switch p.current().tkntype {
		case tknComma:
			p.next()
			break
		case tknValue, tknString:
			t = p.parseTypeDeclaration(t)
			break
		}
	}
	p.enforceNext(tknCloseBracket, "Parameters must close with ')'") // eat ")"
	return t
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
	if !p.isAliasAvailable(alias) {
		p.addError(fmt.Sprintf(errDuplicateAlias, alias, p.prototype.key.key))
	} else {
		p.prototype.fieldAliases[alias] = f
	}
}

func (p *parser) addElementAlias(alias string, e *ProtoElement) {
	if p.prototype.elementAliases == nil {
		p.prototype.elementAliases = make(map[string]*ProtoElement)
	}
	if !p.isAliasAvailable(alias) {
		p.addError(fmt.Sprintf(errDuplicateAlias, alias, p.prototype.key.key))
	} else {
		p.prototype.elementAliases[alias] = e
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
	//	fmt.Println("VALIDATING")
	if p.scope != nil {
		for k, v := range p.prototype.fields {
			//fmt.Printf("key: %s\n", k)
			if v.key.min > len(p.scope.fields[k]) {
				//fmt.Println("MIN")
				p.addError(fmt.Sprintf(errInsufficientFields, v.key.min, k, p.scope.key.key, len(p.scope.fields[k])))
			} else if v.key.max < len(p.scope.fields[k]) {
				//fmt.Println("MAX")
				p.addError(fmt.Sprintf(errDuplicateField, v.key.max, k, p.scope.key.key, len(p.scope.fields[k])))
			}
		}
	}

}

func parseElementClosure(p *parser) {
	p.validateCompleteElement()
	p.prototype = p.prototype.parent
	if p.scope != nil {
		p.scope = p.scope.parent
	}
	p.next()
}

func parseTextAlias(p *parser) {
	p.next() // eat alias
	alias := p.lexer.tokenString(p.next())
	p.next() // eat =
	next := p.next()
	value := strval(p.lexer.tokenString(next))
	p.addTextAlias(alias, TextAlias{value, next.tkntype == tknValue})
}

func (p *parser) addTextAlias(alias string, ta TextAlias) {
	if p.isAliasAvailable(alias) {
		p.prototype.textAliases[alias] = ta
	} else {
		p.addError(errDuplicateAlias)
	}
}

func (p *parser) importValidateConstructs() {
	p.constructs = []construct{
		construct{"field", isField, parseField},
		construct{"element", isElement, parseElement},
		construct{"element closure", isElementClosure, parseElementClosure},
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
	if p.current().tkntype != tknOpenBracket {
		return
	}
	p.next() // eat "("
	// short circuit if no parameters
	if p.current().tkntype == tknCloseBracket {
		return
	}
	e.parameters = make([]*Value, 0)
	e.parameters = p.parseValue(e.parameters)
	for p.current().tkntype == tknComma {
		p.next()
		e.parameters = p.parseValue(e.parameters)
	}
	p.next() // eat ")"'
}

func (p *parser) importPrototypeConstructs() {
	p.constructs = []construct{
		construct{"prototype field alias", isFieldAlias, parseFieldAlias},
		construct{"prototype element alias", isElementAlias, parseElementAlias},
		construct{"text alias", isTextAlias, parseTextAlias},
		construct{"prototype field", isPrototypeField, parsePrototypeField},
		construct{"prototype element", isPrototypeElement, parsePrototypeElement},
		construct{"element closure", isElementClosure, parseElementClosure},
		construct{"discovered alias", isDiscoveredAlias, parseDiscoveredAlias},
	}
}

func (p *parser) end() {
	p.validateCompleteElement()
}

func (p *parser) validateType(validType *TypeDeclaration, fv *Value) bool {
	if validType.isArray {
		//fmt.Println("array")
		if fv.values == nil {
			return false
		}
		//fmt.Printf("fv: %d\n", len(fv.values))
		if validType.max < len(fv.values) && validType.max != 0 {
			//fmt.Println("max")
			return false // p.addError(fmt.Sprintf(errArrayMaximum, key, p.scope.key.key, prototype.key.min, len(fv.values)))
		} else if validType.min > len(fv.values) && validType.min != 0 {
			//fmt.Println("min")
			return false //	p.addError(fmt.Sprintf(errArrayMinimum, key, p.scope.key.key, prototype.key.min, len(fv.values)))
		}
		if validType.types != nil && len(validType.types) != 0 {
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
		} else {

			// should never use below condition?
			if validType.value == nil {
				return false
			}

			return validType.value.MatchString(fv.value)
		}
	} else {
		if validType.types != nil {
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
		} else {
			return validType.value.MatchString(fv.value)
		}
	}
	return true
}

func strval(data string) string {
	cp := make([]byte, len(data))
	copy(cp, []byte(data))
	s := string(cp)
	if strings.HasPrefix(s, "\"") {
		s = strings.TrimPrefix(s, "\"")
		s = strings.TrimSuffix(s, "\"")
	}
	return s
}

func (p *parser) validateField(key string, f *Field) bool {
	//fmt.Printf("in scope %s\n", p.scope.key.key)
	prototype := p.prototype.fields[key]
	if prototype == nil {
		return false
	}
	for _, v := range f.values {
		matched := false
		for _, t := range prototype.types {
			if p.validateType(t, v) {
				matched = true
			}
		}
		if !matched {
			p.addError(fmt.Sprintf(errUnmatchedFieldValue, key))
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
