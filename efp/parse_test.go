package efp

import (
	"fmt"
	"testing"
)

func TestParseSimpleFieldValid(t *testing.T) {
	p := basicParser("name : string")
	assertNow(t, len(p.lexer.tokens) == 3, "wrong token length")
	parsePrototypeField(p)
	assertNow(t, len(p.prototype.fields) == 1, "wrong field length")
	// valid example
	p.createParseString(`name = "ender"`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assertNow(t, len(p.scope.fields["name"][0].value.children) == 1, "wrong children number")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "xxx")

}

func TestParseSimpleFieldInvalid(t *testing.T) {
	p := basicParser("name : string")
	parsePrototypeField(p)
	// invalid example
	p.createParseString(`name = ender`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
}

func TestParseArrayFieldValid(t *testing.T) {
	p := new(parser)
	p.createPrototypeString("name : [string]")
	parsePrototypeField(p)
	p.createParseString(`name = ["ender", "me"]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")
}

func TestParseArrayFieldMinimumValid(t *testing.T) {
	// valid
	p := basicParser("name : [2:string]")
	parsePrototypeField(p)
	p.createParseString(`name = ["ender", "me"]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

}

func TestParseArrayFieldMinimumInvalid(t *testing.T) {
	p := basicParser("name : [2:string]")
	// invalid
	p.createParseString(`name = ["ender", "me", "him"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldMaximumValid(t *testing.T) {
	// valid
	p := basicParser("name : [string:2]")
	p.createParseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")
}

func TestParseArrayFieldMaximumInvalid(t *testing.T) {
	// valid
	p := basicParser("name : [string:2]")
	parsePrototypeField(p)
	// invalid
	p.createParseString(`name = ["ender", "me", "him"]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldFixedValid(t *testing.T) {
	// valid
	p := basicParser("name : [2:string:2]")
	parsePrototypeField(p)
	p.createParseString(`name = ["ender", "me"]`)
	parseField(p)
	fmt.Println(p.errs)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")
}

func TestParseArrayFieldFixedInvalid(t *testing.T) {
	p := basicParser("name : [2:string:2]")
	parsePrototypeField(p)
	// invalid
	p.createParseString(`name = ["ender", "me", "him"]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.createParseString(`name = ["ender"]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldDisjunctionValid(t *testing.T) {
	// valid
	p := basicParser("name : [string|int]")
	parsePrototypeField(p)

	p.createParseString(`name = ["ender", "me"]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.createParseString(`name = [6, 7]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.createParseString(`name = ["ender", 6]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionInvalid(t *testing.T) {
	// valid
	p := basicParser("name : [string|int]")
	parsePrototypeField(p)
	// invalid
	p.createParseString(`name = [true, false]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionMinimumValid(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int]")
	parsePrototypeField(p)

	p.createParseString(`name = ["ender", "me"]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.createParseString(`name = [6, 7]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.createParseString(`name = ["ender", 6]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionMinimumInvalid(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int]")
	parsePrototypeField(p)
	// invalid
	p.createParseString(`name = [true, false]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.createParseString(`name = ["a"]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.createParseString(`name = [6]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionMaximumValid(t *testing.T) {
	// valid
	p := basicParser("name : [string|int:2]")
	parsePrototypeField(p)

	p.createParseString(`name = ["ender", "me"]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.createParseString(`name = [6, 7]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.createParseString(`name = ["ender", 6]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionMaximumInvalid(t *testing.T) {
	// valid
	p := basicParser("name : [string|int:2]")
	parsePrototypeField(p)

	// invalid
	p.createParseString(`name = [false, true]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.createParseString(`name = ["a", "b", "c"]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.createParseString(`name = [6, 7, 8]`)
	parseField(p)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionFixedValid(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int:2]")
	parsePrototypeField(p)

	p.createParseString(`name = ["ender", "me"]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.createParseString(`name = [6, 7]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.createParseString(`name = ["ender", 6]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionFixedInvalid(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int:2]")

	// invalid
	p.createParseString(`name = [false, false]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.createParseString(`name = ["a", "b", "c"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.createParseString(`name = [6, 7, 8]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldTwoDimensionalDisjunction(t *testing.T) {
	// valid
	p := basicParser("name : [2:[2:string|int:2]:2]")
	parsePrototypeField(p)

	p.createParseString(`name = [["ender", "me"], ["me", "ender"]]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
}

func TestParseArrayFieldTwoDimensionalDisjunctionArrays(t *testing.T) {
	// valid
	p := basicParser("name : [2:[2:string:2|[2:int:2]:2]")
	parsePrototypeField(p)
	p.createParseString(`name = [["ender", "me"], ["me", "ender"]]`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
}
