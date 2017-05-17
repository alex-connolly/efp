package efp

import "testing"

func TestPrototypeFieldAlias(t *testing.T) {
	// test only field alias
	p := basicParser(`alias x = key : "value"`)
	assert(t, isFieldAlias(p), "not field alias")
	parseFieldAlias(p)
	assertNow(t, p.prototype.fieldAliases["x"] != nil, "")
	assertNow(t, p.prototype.fieldAliases["x"].key.key == "key", "")
}

func TestPrototypeElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x = key {}`)
	assert(t, isElementAlias(p), "")
	parseElementAlias(p)
	assertNow(t, p.prototype.elementAliases["x"] != nil, "")
	assertNow(t, p.prototype.elementAliases["x"].key.key == "key", "")
}

func TestPrototypeRecursiveElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x = key {}`)
	assert(t, isElementAlias(p), "")
	parseElementAlias(p)
	assertNow(t, p.prototype.elementAliases["x"] != nil, "")
	assertNow(t, p.prototype.elementAliases["x"].key.key == "key", "")
}

func basicParser(data string) *parser {
	p := new(parser)
	p.createPrototypeString(data)
	return p
}

func TestPrototypeFieldBasic(t *testing.T) {
	p := basicParser("name : string")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil && p.prototype.fields["name"] != nil, "")
	assertNow(t, p.prototype.fields["name"].types.value == "string", "")
}

func TestPrototypeFieldBasicDisjunction(t *testing.T) {
	p := basicParser("name : string|int|float")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil && p.prototype.fields["name"] != nil, "")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "wrong length")
	assertNow(t, len(p.prototype.fields["name"][0].value.children) == 3, "wrong value length")
	assertNow(t, p.prototype.fields["name"][0].value.children[0].value == "string", "")
	assertNow(t, p.prototype.fields["name"][0].value.children[1].value == "int", "")
	assertNow(t, p.prototype.fields["name"][0].value.children[2].value == "float", "")
}

func TestPrototypeFieldComplexDisjunction(t *testing.T) {
	p := basicParser(`name : string|"a-zA-Z"|["[abc]{5}":2]`)
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields["name"] != nil, "failed for name")
	assertNow(t, p.prototype.fields["name"][0].value != nil, "failed for value")
	assertNow(t, p.prototype.fields["name"][0].value.children != nil, "failed for nil children")
	assertNow(t, len(p.prototype.fields["name"][0].value.children) == 3, "failed for children length")
	assert(t, p.prototype.fields["name"][0].value.children[0].value == "string", "not string")
	assert(t, p.prototype.fields["name"][0].value.children[1].value == "a-zA-Z", "wrong regex "+p.prototype.fields["name"][0].value.children[1].value)

	assertNow(t, p.prototype.fields["name"][0].value.children[2].children != nil, "children 2 children == nil")
	assert(t, p.prototype.fields["name"][0].value.children[2].children[0].value == "[abc]{5}", "not abc")
	assert(t, p.prototype.fields["name"][0].value.children[2].isArray, "not array")
}

func TestPrototypeFieldAliased(t *testing.T) {

}

func TestPrototypeFieldArray(t *testing.T) {
	p := basicParser("name : [string]")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "standard shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "standard wrong length")
}

func TestPrototypeFieldArrayMinimum(t *testing.T) {
	p := basicParser("name : [2:string]")
	assertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "minimum shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "minimum wrong length")
}

func TestPrototypeFieldArrayMaximum(t *testing.T) {
	p := basicParser("name : [string:2]")
	assertNow(t, len(p.lexer.tokens) == 7, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "maximum shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "maximum wrong length")
}

func TestPrototypeFieldArrayFixed(t *testing.T) {
	p := basicParser("name : [2:string:2]")
	assertNow(t, len(p.lexer.tokens) == 9, "wrong token length")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "fixed shouldn't be nil")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "fixed wrong length")
}

func TestPrototypeFieldRegex(t *testing.T) {
	p := basicParser(`"[a-z]+" : string`)
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.prototype.fields["[a-z]+"]) == 1, "r wrong length")
}

func TestPrototypeFieldRegexEmptyBounds(t *testing.T) {
	p := basicParser(`<"[a-z]+"> : string`)
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.prototype.fields["[a-z]+"]) == 1, "r wrong length")
}

func TestPrototypeFieldRegexMinimumBounds(t *testing.T) {
	p := basicParser(`<2:"[a-z]+"> : string`)
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.prototype.fields["[a-z]+"]) == 1, "r wrong length")
}

func TestPrototypeFieldRegexMaximumBounds(t *testing.T) {
	p := basicParser(`<"[a-z]+":2> : string`)
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil, "r shouldn't be nil")
	assertNow(t, len(p.prototype.fields["[a-z]+"]) == 1, "r wrong length")
}
