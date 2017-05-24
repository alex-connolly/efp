package efp

import "regexp"

// Field ...
type Field struct {
	alias  string
	key    *Key
	values []*Value
}

// Key ...
type Key struct {
	key   string
	regex *regexp.Regexp
	min   int
	max   int
}

// Value ...
type Value struct {
	value  string
	values []*Value
}

// TypeDeclaration ...
type TypeDeclaration struct {
	isArray bool
	types   []*TypeDeclaration
	value   *regexp.Regexp
	min     int
	max     int
}

// ProtoField ...
type ProtoField struct {
	key   *Key
	types []*TypeDeclaration
}

// ProtoElement ...
type ProtoElement struct {
	alias          string
	key            *Key
	parent         *ProtoElement
	parameters     []*TypeDeclaration
	fields         map[string]*ProtoField
	fieldAliases   map[string]*ProtoField
	elements       map[string]*ProtoElement
	elementAliases map[string]*ProtoElement
	aliases        []string
	textAliases    map[string]TextAlias // leave as token for recursion
}

// TextAlias ...
type TextAlias struct {
	value       string
	isRecursive bool
}

// Element ...
type Element struct {
	alias      string
	key        *Key
	parent     *Element
	parameters []*Value
	elements   map[string][]*Element
	fields     map[string][]*Field
}

func (p *ProtoElement) addStandardAliases() {
	p.textAliases = standards
}
