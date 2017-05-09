package efp

import "testing"

func TestParserFieldAlias(t *testing.T) {
	// test only field alias
	p := basicParser(`alias x : key = "value"`)
	if !isFieldAlias(p) {
		t.Fail()
	}
	parseFieldAlias(p)
	if p.prototype.declaredFieldAliases["x"] == nil {
		t.Fail()
	}
	if len(p.prototype.declaredFieldAliases["x"]) != 1 {
		t.Fail()
	}
	if p.prototype.declaredFieldAliases["x"][0].key != "key" {
		t.Fail()
	}
}

func TestParserElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x : key {}`)
	if !isElementAlias(p) {
		t.Fail()
	}
	parseElementAlias(p)
	if p.prototype.declaredElementAliases["x"] == nil {
		t.Fail()
	}
	if p.prototype.declaredElementAliases["x"][0].key != "key" {
		t.Fail()
	}
}

func TestParserRecursiveElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x : key {}`)
	if !isElementAlias(p) {
		t.Fail()
	}
	parseElementAlias(p)
	if p.prototype.declaredElementAliases["x"] == nil {
		t.Fail()
	}
	if p.prototype.declaredElementAliases["x"][0].key != "key" {
		t.Fail()
	}
}

func TestParserIdentifierMethods(t *testing.T) {
	p := basicParser("alias x : a = 4")
	if !isFieldAlias(p) {
		t.Fail()
	}
	p = basicParser("alias x : element {}")
	if !isElementAlias(p) {
		t.Fail()
	}
	p = basicParser("element {}")
	if !isElement(p) {
		t.Fail()
	}

	p = basicParser(`x = 5`)
	if !isField(p) {
		t.Fail()
	}
	p = basicParser(`}`)
	if !isElementClosure(p) {
		t.Fail()
	}
	p = basicParser(`element("name", 5){}`)
	if !isElement(p) {
		t.Fail()
	}
}

func basicParser(data string) *parser {
	p := new(parser)
	p.lexer = lex([]byte(data))
	p.prototype = new(element)
	p.index = 0
	p.importPrototypeConstructs()
	return p
}

func TestParserPrototypeFieldBasic(t *testing.T) {
	p := basicParser("name : string")
	if !isPrototypeField(p) {
		t.Fail()
	}
	parsePrototypeField(p)
	if p.prototype.fields == nil || p.prototype.fields["name"] == nil {
		t.Fail()
	}
	if len(p.prototype.fields["name"]) != 1 {
		t.Fail()
	}
}

func TestParserPrototypeFieldBasicDisjunction(t *testing.T) {
	p := basicParser("name : string|int|float")
	if !isPrototypeField(p) {
		t.Fail()
	}

}

func TestParserPrototypeFieldComplexDisjunction(t *testing.T) {
	p := basicParser(`name : string|"a-zA-Z"|["[abc]{5}":2]`)
	if !isPrototypeField(p) {
		t.Fail()
	}
}

func TestParserPrototypeFieldAliased(t *testing.T) {

}

func TestParserPrototypeFieldArray(t *testing.T) {
	p := basicParser("name : [string]")
	if !isPrototypeField(p) {
		t.Fail()
	}
	parsePrototypeField(p)
	if p.prototype.fields == nil || p.prototype.fields["name"] == nil {
		t.Fail()
	}
	if len(p.prototype.fields["name"]) != 1 {
		t.Fail()
	}
	p = basicParser("name : [2:string]")
	if !isPrototypeField(p) {
		t.Fail()
	}
	parsePrototypeField(p)
	if p.prototype.fields == nil || p.prototype.fields["name"] == nil {
		t.Fail()
	}
	if len(p.prototype.fields["name"]) != 1 {
		t.Fail()
	}
	p = basicParser("name : [string:2]")
	if !isPrototypeField(p) {
		t.Fail()
	}
	parsePrototypeField(p)
	if p.prototype.fields == nil || p.prototype.fields["name"] == nil {
		t.Fail()
	}
	if len(p.prototype.fields["name"]) != 1 {
		t.Fail()
	}
	p = basicParser("name : [2:string:2]")
	if !isPrototypeField(p) {
		t.Fail()
	}
	parsePrototypeField(p)
	if p.prototype.fields == nil || p.prototype.fields["name"] == nil {
		t.Fail()
	}
	if len(p.prototype.fields["name"]) != 1 {
		t.Fail()
	}
}
