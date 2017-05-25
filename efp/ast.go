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

// Type ...
func (f *ProtoField) Type(indices ...int) *TypeDeclaration {
	current := f.types
	for i := 0; i < len(indices)-1; i++ {
		if current == nil {
			return nil
		}
		current = current[indices[i]].types

	}
	return current[indices[len(indices)-1]] // trim in case
}

// TypeValue ...
func (f *ProtoField) TypeValue(indices ...int) string {
	t := f.Type(indices...)
	if t == nil {
		return ""
	}
	return t.value.String()
}

// Field ...
func (f *ProtoElement) Field(name string) *ProtoField {
	return f.fields[name]
}

// Types ...
func (f *ProtoField) Types(indices ...int) []*TypeDeclaration {
	current := f.types
	for i := 0; i < len(indices); i++ {
		current = current[indices[i]].types
	}
	return current
}

// Value ...
func (f *Field) Value(indices ...int) string {
	current := f.values
	for i := 0; i < len(indices)-1; i++ {
		if current == nil {
			return ""
		}
		current = current[indices[i]].values
	}
	return strval(current[indices[len(indices)-1]].value) // trim in case
}

// Values ...
func (f *Field) Values(indices ...int) []*Value {
	current := f.values
	for i := 0; i < len(indices); i++ {
		current = current[indices[i]].values
	}
	return current
}

// Element ...
func (e *Element) Element(name string, index int) *Element {
	return e.elements[name][index]
}

// Field ...
func (e *Element) Field(name string, index int) *Field {
	return e.fields[name][index]
}

func (e *ProtoElement) addStandardAliases() {
	// reassign everything to keep standards as a separate map
	// better for testing
	e.textAliases = make(map[string]TextAlias)
	for k, v := range standards {
		e.textAliases[k] = v
	}
}
