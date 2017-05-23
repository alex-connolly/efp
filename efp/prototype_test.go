package efp

import (
	"fmt"
	"testing"
)

func TestPrototypeFieldAlias(t *testing.T) {
	// test only field alias
	p := createPrototypeParserString(`alias x = key : "value"`)
	assert(t, isFieldAlias(p), "not field alias")
	parseFieldAlias(p)
	assertNow(t, p.prototype.fieldAliases["x"] != nil, "")
	assertNow(t, p.prototype.fieldAliases["x"].key.key == "key", "")
}

func TestPrototypeElementAlias(t *testing.T) {
	// test only element alias
	p, _ := PrototypeString(`alias x = key {}`)
	assertNow(t, p.elementAliases["x"] != nil, "x is nil")
	assertNow(t, p.elementAliases["x"].key.key == "key",
		fmt.Sprintf("key should be 'key', not %s", p.elementAliases["x"].key.key))
}

func TestPrototypeRecursiveElementAlias(t *testing.T) {
	/* test only element alias
	p := createPrototypeParserString(`alias x = key {}`)
	assert(t, isElementAlias(p), "not element alias")
	parseElementAlias(p)
	assertNow(t, p.prototype.elementAliases["x"] != nil, "x is nil")
	assertNow(t, p.prototype.elementAliases["x"].key.key == "key",
		fmt.Sprintf("key should be 'key', not %s", p.prototype.elementAliases["x"].key.key))*/
}

func TestPrototypeFieldBasic(t *testing.T) {
	p := createPrototypeParserString("name : string")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype != nil, "prototype is nil")
	assertNow(t, p.prototype.fields != nil && p.prototype.fields["name"] != nil, "fields is nil")
	assertNow(t, len(p.prototype.fields["name"].types) == 1, "wrong type length")
	assertNow(t, p.prototype.fields["name"].types[0].value.String() == standards["string"].value, "wrong type")
}

func TestPrototypeFieldBasicDisjunction(t *testing.T) {
	p := createPrototypeParserString("name : string|int|float")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil && p.prototype.fields["name"] != nil, "")
	assertNow(t, len(p.prototype.fields["name"].types) == 3, "wrong length")
	assertNow(t, p.prototype.fields["name"].types[0].value.String() == standards["string"].value, "")
	assertNow(t, p.prototype.fields["name"].types[1].value.String() == standards["int"].value, "")
	assertNow(t, p.prototype.fields["name"].types[2].value.String() == standards["float"].value, "")
}

func TestPrototypeFieldComplexDisjunction(t *testing.T) {
	p := createPrototypeParserString(`name : string|"[a-zA-Z]+"|["[abc]{5}":2]`)
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields["name"] != nil, "failed for name")
	assertNow(t, p.prototype.fields["name"].types != nil, "failed for nil children")
	assertNow(t, len(p.prototype.fields["name"].types) == 3, "failed for children length")
	assert(t, p.prototype.fields["name"].types[1].value.String() == "[a-zA-Z]+", "wrong regex "+p.prototype.fields["name"].types[1].value.String())

	assertNow(t, p.prototype.fields["name"].types[2].types != nil, "children 2 children == nil")
	assert(t, p.prototype.fields["name"].types[2].types[0].value.String() == "[abc]{5}", "not abc")
	assert(t, p.prototype.fields["name"].types[2].isArray, "not array")
}

func TestPrototypeFieldAliased(t *testing.T) {

}

func TestPrototypeFieldArray(t *testing.T) {
	p := createPrototypeParserString("name : [string]")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "standard shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"].types) == 1, "standard wrong length")
}

func TestPrototypeFieldArrayMinimum(t *testing.T) {
	p := createPrototypeParserString("name : [2:string]")
	assertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "minimum shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"].types) == 1, "minimum wrong length")
}

func TestPrototypeFieldArrayMaximum(t *testing.T) {
	p := createPrototypeParserString("name : [string:2]")
	assertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "maximum shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"].types) == 1, "maximum wrong length")
}

func TestPrototypeFieldArrayFixed(t *testing.T) {
	p := createPrototypeParserString("name : [2:string:2]")
	assertNow(t, len(p.lexer.tokens) == 9, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "fixed shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"].types) == 1, "fixed wrong length")
}

func TestPrototypeFieldRegex(t *testing.T) {
	p, errs := PrototypeString(`"[a-z]+" : string`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, "r wrong length")
}

func TestPrototypeFieldRegexEmptyBounds(t *testing.T) {
	p, errs := PrototypeString(`<"[a-z]+"> : string`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.fields) == 1, "wrong field length")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, "r wrong length")
}

func TestPrototypeFieldRegexMinimumBounds(t *testing.T) {
	p, errs := PrototypeString(`<2:"[a-z]+"> : string`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.fields) == 1, "wrong field length")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, "r wrong length")
}

func TestPrototypeFieldRegexMaximumBounds(t *testing.T) {
	p, errs := PrototypeString(`<"[a-z]+":2> : string`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.fields) == 1, "wrong field length")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, fmt.Sprintf("r wrong length (%d)", len(p.fields["[a-z]+"].types)))
}
