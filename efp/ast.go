package efp

import "regexp"

type Field struct {
	alias  string
	key    *Key
	values []*Value
}

type Key struct {
	key   string
	regex *regexp.Regexp
	min   int
	max   int
}

type Value struct {
	value  string
	values []*Value
}

type TypeDeclaration struct {
	isArray bool
	types   []*TypeDeclaration
	value   *regexp.Regexp
	min     int
	max     int
}
type ProtoField struct {
	key   *Key
	types []*TypeDeclaration
}

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

type TextAlias struct {
	value       string
	isRecursive bool
}

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
