package efp

import "testing"

func TestFieldAlias(t *testing.T) {

}

func TestElementAlias(t *testing.T) {

}

func TestRecursiveElementAlias(t *testing.T) {

}

func basicParser(data string) *parser {
	p := new(parser)
	p.lexer = lex([]byte(data))
	p.prototype = new(element)
	p.offset = 0
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
	if p.prototype.fields["name"][0].value[0] != "string" {
		t.Fail()
	}
}

func TestPrototypeFieldBasicDisjunction(t *testing.T) {
	p := basicParser("name : string|int|float")
	if !isPrototypeField(p) {
		t.Fail()
	}
	if p.prototype.fields["name"][0].value[1] != "int" {

	}
}

func TestPrototypeFieldComplexDisjunction(t *testing.T) {
	p := basicParser(`name : string|"a-zA-Z"|["[abc]{5}":2]`)
	if !isPrototypeField(p) {
		t.Fail()
	}
	if p.prototype.fields["name"][0].value[1] != "int" {

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
	if p.prototype.fields["name"][0].value[0] != "string" {
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
	if p.prototype.fields["name"][0].value[0] != "string" {
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
	if p.prototype.fields["name"][0].value[0] != "string" {
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
	if p.prototype.fields["name"][0].value[0] != "string" {
		t.Fail()
	}
}
