package efp

import "testing"

func TestPrototypeFieldAlias(t *testing.T) {
	// test only field alias
	p := basicParser(`alias x : key = "value"`)
	assert(t, isFieldAlias(p), "not field alias")
	parseFieldAlias(p)
	assertNow(t, p.prototype.declaredFieldAliases["x"] != nil, "")
	assertNow(t, len(p.prototype.declaredFieldAliases["x"]) == 1, "")
	assertNow(t, p.prototype.declaredFieldAliases["x"][0].key == "key", "")
}

func TestPrototypeElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x : key {}`)
	assert(t, isElementAlias(p), "")
	parseElementAlias(p)
	assertNow(t, p.prototype.declaredElementAliases["x"] != nil, "")
	assertNow(t, p.prototype.declaredElementAliases["x"][0].key == "key", "")
}

func TestPrototypeRecursiveElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x : key {}`)
	assert(t, isElementAlias(p), "")
	parseElementAlias(p)
	assertNow(t, p.prototype.declaredElementAliases["x"] != nil, "")
	assertNow(t, p.prototype.declaredElementAliases["x"][0].key == "key", "")
}

func basicParser(data string) *parser {
	p := new(parser)
	p.prototypeString(data)
	return p
}

func TestPrototypeFieldBasic(t *testing.T) {
	p := basicParser("name : string")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil && p.prototype.fields["name"] != nil, "")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "wrong length")
	assertNow(t, p.prototype.fields["name"][0].value.children[0].value == "string", "")
}

func TestPrototypeFieldBasicDisjunction(t *testing.T) {
	p := basicParser("name : string|int|float")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields != nil && p.prototype.fields["name"] != nil, "")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "wrong length")
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
	assert(t, p.prototype.fields["name"][0].value.children[0].value == "string", "")
	assert(t, p.prototype.fields["name"][0].value.children[1].value == "a-zA-Z", "")
	assert(t, p.prototype.fields["name"][0].value.children[2].value == "[abc]{5}", "")
	assert(t, p.prototype.fields["name"][0].value.children[2].isArray, "not array")
	assert(t, p.prototype.fields["name"][0].value.children[2].value == "[abc]{5}", "")
}

func TestPrototypeFieldAliased(t *testing.T) {

}

func TestPrototypeFieldArray(t *testing.T) {
	p := basicParser("name : [string]")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields == nil || p.prototype.fields["name"] == nil, "")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "")

	p = basicParser("name : [2:string]")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields == nil || p.prototype.fields["name"] == nil, "")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "")

	p = basicParser("name : [string:2]")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields == nil || p.prototype.fields["name"] == nil, "")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "")

	p = basicParser("name : [2:string:2]")
	assert(t, isPrototypeField(p), "not prototype field")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields == nil || p.prototype.fields["name"] == nil, "")
	assertNow(t, len(p.prototype.fields["name"]) == 1, "")

}
