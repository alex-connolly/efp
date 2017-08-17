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
	valueAliases   map[string]ValueAlias // leave as token for recursion
}

// ValueAlias ...
type ValueAlias struct {
	values      []*TypeDeclaration
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
func (e *ProtoElement) Field(name string) *ProtoField {
	return e.fields[name]
}

// Element ...
func (e *ProtoElement) Element(name string) *ProtoElement {
	return e.elements[name]
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
	// Value() shouldn't actually be called (makes no sense), but will be:
	if len(indices) == 0 {
		return f.Value(0)
	}

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

// Value ...
func (v *Value) Value() string {
	return strval(v.value)
}

// Parameter ...
func (e *Element) Parameter(index int) *Value {
	return e.parameters[index]
}

// Parameters ...
func (e *Element) Parameters() []*Value {
	return e.parameters
}

// FirstElement ...
func (e *Element) FirstElement(name string) *Element {
	return e.elements[name][0]
}

// Element ...
func (e *Element) Element(name string, index int) *Element {
	return e.elements[name][index]
}

// Elements ...
func (e *Element) Elements(name string) []*Element {
	return e.elements[name]
}

// Fields ...
func (e *Element) Fields(name string) []*Field {
	return e.fields[name]
}

// Field ...
func (e *Element) Field(name string, index int) *Field {
	return e.fields[name][index]
}

// FirstField ...
func (e *Element) FirstField(name string) *Field {
	return e.fields[name][0]
}
