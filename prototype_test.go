package efp

import (
	"fmt"
	"testing"

	"github.com/end-r/goutil"
)

func TestPrototypeFieldAlias(t *testing.T) {
	// test only field alias
	p, errs := PrototypeString(`alias x = key : int`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fieldAliases["x"] != nil, "x is nil")
	goutil.AssertNow(t, p.fieldAliases["x"].key.key == "key", "wrong key for x")
	goutil.AssertNow(t, p.fieldAliases["x"].TypeValue(0) == standards["int"], "wrong type for field")
}

func TestPrototypeElementAlias(t *testing.T) {
	// test only element alias
	p, _ := PrototypeString(`alias x = key {}`)
	goutil.AssertNow(t, p.elementAliases["x"] != nil, "x is nil")
	goutil.AssertNow(t, p.elementAliases["x"].key.key == "key",
		fmt.Sprintf("key should be 'key', not %s", p.elementAliases["x"].key.key))
}

func TestPrototypeRecursiveElementAlias(t *testing.T) {
	/* test only element alias
	p := createPrototypeParserString(`alias x = key {}`)
	goutil.Assert(t, isElementAlias(p), "not element alias")
	parseElementAlias(p)
	goutil.AssertNow(t, p.prototype.elementAliases["x"] != nil, "x is nil")
	goutil.AssertNow(t, p.prototype.elementAliases["x"].key.key == "key",
		fmt.Sprintf("key should be 'key', not %s", p.prototype.elementAliases["x"].key.key))*/
}

func TestPrototypeFieldBasic(t *testing.T) {
	p := createPrototypeParserString("name : string")
	goutil.Assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	goutil.AssertNow(t, p.prototype != nil, "prototype is nil")
	goutil.AssertNow(t, p.prototype.fields != nil && p.prototype.Field("name") != nil, "fields is nil")
	goutil.AssertNow(t, len(p.prototype.Field("name").types) == 1, "wrong type length")
	goutil.AssertNow(t, p.prototype.Field("name").TypeValue(0) == standards["string"], "wrong type")
}

func TestPrototypeFieldBasicDisjunction(t *testing.T) {
	p := createPrototypeParserString("name : string|int|float")
	goutil.Assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	goutil.AssertNow(t, p.prototype.fields != nil && p.prototype.Field("name") != nil, "")
	goutil.AssertNow(t, len(p.prototype.Field("name").types) == 3, "wrong length")
	goutil.AssertNow(t, p.prototype.Field("name").TypeValue(0) == standards["string"], "")
	goutil.AssertNow(t, p.prototype.Field("name").TypeValue(1) == standards["int"], "")
	goutil.AssertNow(t, p.prototype.Field("name").TypeValue(2) == standards["float"], "")
}

func TestPrototypeFieldComplexDisjunction(t *testing.T) {
	p, errs := PrototypeString(`name : string|"[a-zA-Z]+"|["[abc]{5}":2]`)
	goutil.AssertNow(t, errs == nil, "errs must be nil")
	goutil.AssertNow(t, p.Field("name") != nil, "failed for name")
	goutil.AssertNow(t, p.Field("name").types != nil, "failed for nil children")

	goutil.AssertNow(t, len(p.Field("name").types) == 3, "failed for children length")
	goutil.Assert(t, p.Field("name").TypeValue(1) == "[a-zA-Z]+", "wrong regex "+p.Field("name").TypeValue(1))

	goutil.AssertNow(t, len(p.Field("name").types[2].types) == 1, "children 2 children length incorrect")
	goutil.Assert(t, p.Field("name").TypeValue(2, 0) == "[abc]{5}", "not abc")
	//goutil.Assert(t, p.Field("name").TypeValue(0).isArray, "not array")
}

func TestPrototypeFieldAliased(t *testing.T) {

}

func TestPrototypeFieldArray(t *testing.T) {
	p := createPrototypeParserString("name : [string]")
	goutil.Assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	goutil.AssertNow(t, p.prototype.fields != nil, "standard shouldn't be nil")
	goutil.AssertNow(t, len(p.prototype.Field("name").types) == 1, "standard wrong length")
}

func TestPrototypeFieldArrayMinimum(t *testing.T) {
	p := createPrototypeParserString("name : [2:string]")
	goutil.AssertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	goutil.Assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	goutil.AssertNow(t, p.prototype.fields != nil, "minimum shouldn't be nil")
	goutil.AssertNow(t, len(p.prototype.Field("name").types) == 1, "minimum wrong length")
}

func TestPrototypeFieldArrayMaximum(t *testing.T) {
	p := createPrototypeParserString("name : [string:2]")
	goutil.AssertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	goutil.Assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	goutil.AssertNow(t, p.prototype.fields != nil, "maximum shouldn't be nil")
	goutil.AssertNow(t, len(p.prototype.Field("name").types) == 1, "maximum wrong length")
}

func TestPrototypeFieldArrayFixed(t *testing.T) {
	p := createPrototypeParserString("name : [2:string:2]")
	goutil.AssertNow(t, len(p.lexer.tokens) == 9, "wrong token length")
	goutil.Assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	goutil.AssertNow(t, p.prototype.fields != nil, "fixed shouldn't be nil")
	goutil.AssertNow(t, len(p.prototype.Field("name").types) == 1, "fixed wrong length")
}

func TestPrototypeFieldRegex(t *testing.T) {
	p, errs := PrototypeString(`"[a-z]+" : string`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "r shouldn't be nil")
	goutil.AssertNow(t, len(p.fields["[a-z]+"].types) == 1, "r wrong length")
}

func TestPrototypeFieldRegexEmptyBounds(t *testing.T) {
	p, errs := PrototypeString(`<"[a-z]+"> : string`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "r shouldn't be nil")
	goutil.AssertNow(t, len(p.fields) == 1, "wrong field length")
	goutil.AssertNow(t, len(p.fields["[a-z]+"].types) == 1, "r wrong length")
}

func TestPrototypeFieldRegexMinimumBounds(t *testing.T) {
	p, errs := PrototypeString(`<2:"[a-z]+"> : string`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "r shouldn't be nil")
	goutil.AssertNow(t, len(p.fields) == 1, "wrong field length")
	goutil.AssertNow(t, p.fields["[a-z]+"] != nil, "types must not be nil")
	goutil.AssertNow(t, len(p.fields["[a-z]+"].types) == 1, "r wrong length")
	/*	p, errs = PrototypeString(`alias MIN = 2 <MIN:"[a-z]+"> : string`)
		goutil.Assert(t, errs == nil, "errs should be nil")
		goutil.AssertNow(t, p.fields != nil, "r shouldn't be nil")
		goutil.AssertNow(t, len(p.fields) == 1, "wrong field length")
		goutil.AssertNow(t, p.fields["[a-z]+"] != nil, "field key is nil entry")
		goutil.AssertNow(t, len(p.fields["[a-z]+"].types) == 1, fmt.Sprintf("r wrong length (%d)", len(p.fields["[a-z]+"].types)))*/
}

func TestPrototypeFieldRegexMaximumBounds(t *testing.T) {
	p, errs := PrototypeString(`<"[a-z]+":2> : string`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "r shouldn't be nil")
	goutil.AssertNow(t, len(p.fields) == 1, "wrong field length")
	goutil.AssertNow(t, len(p.Field("[a-z]+").types) == 1, fmt.Sprintf("r wrong length (%d)", len(p.fields["[a-z]+"].types)))
	/*p, errs = PrototypeString(`alias MAX = 2 <"[a-z]+":MAX> : string`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "r shouldn't be nil")
	goutil.AssertNow(t, len(p.fields) == 1, "wrong field length")
	goutil.AssertNow(t, p.fields["[a-z]+"] != nil, "field key is nil entry")
	goutil.AssertNow(t, p.fields["[a-z]+"].types != nil, "types is nil")
	goutil.AssertNow(t, len(p.fields["[a-z]+"].types) == 1, fmt.Sprintf("r wrong length (%d)", len(p.fields["[a-z]+"].types)))*/
}

func TestTwoDimensionalArraySimple(t *testing.T) {
	p, errs := PrototypeString("name : [[string]]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "fields shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")                           // top level array
	goutil.AssertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                        // second level array
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"], "type should be string") //
}

func TestTwoDimensionalArrayDisjunction(t *testing.T) {
	p, errs := PrototypeString("name : [[string|int]]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "fields shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")                           // top level array
	goutil.AssertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                        // second level array
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"], "type should be string") //
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 1) == standards["int"], "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunction(t *testing.T) {
	p, errs := PrototypeString("name : [[string]|[int]]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "fields shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")                           // top level array
	goutil.AssertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                        // second level array
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"], "type should be string") //
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"], "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionMinimum(t *testing.T) {
	p, errs := PrototypeString("name : [2:[string]|[int]]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "fields shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Type(0).min == 2, "wrong minimum value")                            // top level array
	goutil.AssertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                        // second level array
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"], "type should be string") //
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"], "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionMaximum(t *testing.T) {
	p, errs := PrototypeString("name : [[string]|[int]:2]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "fields shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Type(0).max == 2, "wrong maximum value")                            // top level array
	goutil.AssertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                        // second level array
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"], "type should be string") //
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"], "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionFixed(t *testing.T) {
	p, errs := PrototypeString("name : [2:[string]|[int]:2]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "fields shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Type(0).max == 2, "wrong maximum value")
	goutil.AssertNow(t, p.Field("name").Type(0).min == 2, "wrong minimum value")                            // top level array
	goutil.AssertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")                        // second level array
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"], "type should be string") //
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"], "type should be int")
}

func TestTwoDimensionalArrayArrayDisjunctionFixedComplex(t *testing.T) {
	p, errs := PrototypeString("name : [2:[3:string|float:3]|[4:int:4]:2]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields != nil, "fields shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Types() != nil, "types shouldn't be nil")
	goutil.AssertNow(t, p.Field("name").Type(0).max == 2, "wrong maximum value")
	goutil.AssertNow(t, p.Field("name").Type(0).min == 2, "wrong minimum value") // top level array
	goutil.AssertNow(t, p.Field("name").Types(0) != nil, "types 0 shouldn't be nil")

	goutil.Assert(t, p.Field("name").Type(0, 0).max == 3, "wrong maximum value")
	goutil.Assert(t, p.Field("name").Type(0, 0).min == 3, "wrong minimum value")
	goutil.Assert(t, p.Field("name").Type(0, 1).max == 4, "wrong maximum value")
	goutil.Assert(t, p.Field("name").Type(0, 1).min == 4, "wrong minimum value")

	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 0) == standards["string"], "type should be string")
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 0, 1) == standards["float"], "type should be string")
	goutil.AssertNow(t, p.Field("name").TypeValue(0, 1, 0) == standards["int"], "type should be int")
}
