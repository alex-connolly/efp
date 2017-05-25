package efp

import (
	"fmt"
	"testing"
)

func TestPrototypeFieldAlias(t *testing.T) {
	// test only field alias
	p, errs := PrototypeString(`alias x = key : int`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fieldAliases["x"] != nil, "x is nil")
	assertNow(t, p.fieldAliases["x"].key.key == "key", "wrong key for x")
	assertNow(t, p.fieldAliases["x"].TypeValue(0) == standards["int"].value, "wrong type for field")
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
	assertNow(t, p.prototype.fields != nil && p.prototype.Field("name") != nil, "fields is nil")
	assertNow(t, len(p.prototype.Field("name").types) == 1, "wrong type length")
	assertNow(t, p.prototype.Field("name").TypeValue(0) == standards["string"].value, "wrong type")
}

func TestPrototypeFieldBasicDisjunction(t *testing.T) {
	p := createPrototypeParserString("name : string|int|float")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil && p.prototype.Field("name") != nil, "")
	assertNow(t, len(p.prototype.Field("name").types) == 3, "wrong length")
	assertNow(t, p.prototype.Field("name").TypeValue(0) == standards["string"].value, "")
	assertNow(t, p.prototype.Field("name").TypeValue(1) == standards["int"].value, "")
	assertNow(t, p.prototype.Field("name").TypeValue(2) == standards["float"].value, "")
}

func TestPrototypeFieldComplexDisjunction(t *testing.T) {
	p, errs := PrototypeString(`name : string|"[a-zA-Z]+"|["[abc]{5}":2]`)
	assertNow(t, errs == nil, "errs must be nil")
	assertNow(t, p.Field("name") != nil, "failed for name")
	assertNow(t, p.Field("name").types != nil, "failed for nil children")

	assertNow(t, len(p.Field("name").types) == 3, "failed for children length")
	assert(t, p.Field("name").TypeValue(1) == "[a-zA-Z]+", "wrong regex "+p.Field("name").TypeValue(1))

	assertNow(t, len(p.Field("name").types[2].types) == 1, "children 2 children length incorrect")
	assert(t, p.Field("name").TypeValue(2, 0) == "[abc]{5}", "not abc")
	//assert(t, p.Field("name").TypeValue(0).isArray, "not array")
}

func TestPrototypeFieldAliased(t *testing.T) {

}

func TestPrototypeFieldArray(t *testing.T) {
	p := createPrototypeParserString("name : [string]")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "standard shouldn't be nil")
	assertNow(t, len(p.prototype.Field("name").types) == 1, "standard wrong length")
}

func TestPrototypeFieldArrayMinimum(t *testing.T) {
	p := createPrototypeParserString("name : [2:string]")
	assertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "minimum shouldn't be nil")
	assertNow(t, len(p.prototype.Field("name").types) == 1, "minimum wrong length")
}

func TestPrototypeFieldArrayMaximum(t *testing.T) {
	p := createPrototypeParserString("name : [string:2]")
	assertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "maximum shouldn't be nil")
	assertNow(t, len(p.prototype.Field("name").types) == 1, "maximum wrong length")
}

func TestPrototypeFieldArrayFixed(t *testing.T) {
	p := createPrototypeParserString("name : [2:string:2]")
	assertNow(t, len(p.lexer.tokens) == 9, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "fixed shouldn't be nil")
	assertNow(t, len(p.prototype.Field("name").types) == 1, "fixed wrong length")
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
	assertNow(t, p.fields["[a-z]+"] != nil, "types must not be nil")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, "r wrong length")
	p, errs = PrototypeString(`alias MIN = 2 <MIN:"[a-z]+"> : string`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.fields) == 1, "wrong field length")
	assertNow(t, p.fields["[a-z]+"] != nil, "field key is nil entry")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, fmt.Sprintf("r wrong length (%d)", len(p.fields["[a-z]+"].types)))
}

func TestPrototypeFieldRegexMaximumBounds(t *testing.T) {
	p, errs := PrototypeString(`<"[a-z]+":2> : string`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.fields) == 1, "wrong field length")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, fmt.Sprintf("r wrong length (%d)", len(p.fields["[a-z]+"].types)))
	p, errs = PrototypeString(`alias MAX = 2 <"[a-z]+":MAX> : string`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.fields) == 1, "wrong field length")
	assertNow(t, p.fields["[a-z]+"] != nil, "field key is nil entry")
	assertNow(t, p.fields["[a-z]+"].types != nil, "types is nil")
	assertNow(t, len(p.fields["[a-z]+"].types) == 1, fmt.Sprintf("r wrong length (%d)", len(p.fields["[a-z]+"].types)))
}

func TestTwoDimensionalArraySimple(t *testing.T) {
	p, errs := PrototypeString("name : [[string]]")
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "fields shouldn't be nil")
	assertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")                                 // top level array
	assertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                              // second level array
	assertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"].value, "type should be string") //
}

func TestTwoDimensionalArrayDisjunction(t *testing.T) {
	p, errs := PrototypeString("name : [[string|int]]")
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "fields shouldn't be nil")
	assertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")                                 // top level array
	assertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                              // second level array
	assertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"].value, "type should be string") //
	assertNow(t, p.Field("name").TypeValue(0, 0, 1) == standards["int"].value, "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunction(t *testing.T) {
	p, errs := PrototypeString("name : [[string]|[int]]")
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "fields shouldn't be nil")
	assertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")                                 // top level array
	assertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                              // second level array
	assertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"].value, "type should be string") //
	assertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"].value, "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionMinimum(t *testing.T) {
	p, errs := PrototypeString("name : [2:[string]|[int]]")
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "fields shouldn't be nil")
	assertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	assertNow(t, p.Field("name").Type(0).min == 2, "wrong minimum value")                                  // top level array
	assertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                              // second level array
	assertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"].value, "type should be string") //
	assertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"].value, "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionMaximum(t *testing.T) {
	p, errs := PrototypeString("name : [[string]|[int]:2]")
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "fields shouldn't be nil")
	assertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	assertNow(t, p.Field("name").Type(0).max == 2, "wrong maximum value")                                  // top level array
	assertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                              // second level array
	assertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"].value, "type should be string") //
	assertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"].value, "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionFixed(t *testing.T) {
	p, errs := PrototypeString("name : [2:[string]|[int]:2]")
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "fields shouldn't be nil")
	assertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	assertNow(t, p.Field("name").Type(0).max == 2, "wrong maximum value")
	assertNow(t, p.Field("name").Type(0).min == 2, "wrong minimum value")                                  // top level array
	assertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                              // second level array
	assertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"].value, "type should be string") //
	assertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"].value, "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionFixedComplex(t *testing.T) {
	p, errs := PrototypeString("name : [2:[3:string|float:3]|[4:int:4]:2]")
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields != nil, "fields shouldn't be nil")
	assertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	assertNow(t, p.Field("name").Type(0).max == 2, "wrong maximum value")
	assertNow(t, p.Field("name").Type(0).min == 2, "wrong minimum value") // top level array
	assertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")

	assert(t, p.Field("name").Type(0, 0).max == 3, "wrong maximum value")
	assert(t, p.Field("name").Type(0, 0).min == 3, "wrong minimum value")
	assert(t, p.Field("name").Type(0, 1).max == 4, "wrong maximum value")
	assert(t, p.Field("name").Type(0, 1).min == 4, "wrong minimum value")

	assertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"].value, "type should be string")
	assertNow(t, p.Field("name").TypeValue(0, 0, 1) == standards["float"].value, "type should be string")
	assertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"].value, "type should be int")
}
