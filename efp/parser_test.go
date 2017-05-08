package efp

import "testing"

func TestFieldAlias(t *testing.T) {
	// test only field alias
	p := basicParser(`alias x : key = "value"`)
	if !isFieldAlias(p) {
		t.Fail()
	}
	parseFieldAlias(p)
	if p.scope.declaredFieldAliases["x"] == nil {
		t.Fail()
	}
	if p.scope.declaredFieldAliases["x"][0].key != "key" {
		t.Fail()
	}
}

func TestElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x : key {}`)
	if !isElementAlias(p) {
		t.Fail()
	}
	parseElementAlias(p)
	if p.scope.declaredElementAliases["x"] == nil {
		t.Fail()
	}
	if p.scope.declaredElementAliases["x"][0].key != "key" {
		t.Fail()
	}
}

func TestRecursiveElementAlias(t *testing.T) {
	// test only element alias
	p := basicParser(`alias x : key {}`)
	if !isElementAlias(p) {
		t.Fail()
	}
	parseElementAlias(p)
	if p.scope.declaredElementAliases["x"] == nil {
		t.Fail()
	}
	if p.scope.declaredElementAliases["x"][0].key != "key" {
		t.Fail()
	}
}

func TestIdentiferMethods(t *testing.T) {
	p := basicParser("alias x : a = 4")
	if !isFieldAlias(p) {
		t.Fail()
	}
	p = basicParser("alias x : element {}")
	if !isElementAlias(p) {
		t.Fail()
	}
	p = basicParser("alias x : 5")
	if !isTextAlias(p) {
		t.Fail()
	}
	p = basicParser("element {}")
	if !isElement(p) {
		t.Fail()
	}
	p = basicParser(`element("name", 5){}`)
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
}

func basicParser(data string) *parser {
	p := new(parser)
	p.lexer = lex([]byte(data))
	p.prototype = new(element)
	p.index = 0
	p.importPrototypeConstructs()
	return p
}

func TestPrototypeFieldBasic(t *testing.T) {
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

func TestPrototypeFieldBasicDisjunction(t *testing.T) {
	p := basicParser("name : string|int|float")
	if !isPrototypeField(p) {
		t.Fail()
	}

}

func TestPrototypeFieldComplexDisjunction(t *testing.T) {
	p := basicParser(`name : string|"a-zA-Z"|["[abc]{5}":2]`)
	if !isPrototypeField(p) {
		t.Fail()
	}
}

func TestPrototypeFieldAliased(t *testing.T) {

}

func TestPrototypeFieldArray(t *testing.T) {
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
