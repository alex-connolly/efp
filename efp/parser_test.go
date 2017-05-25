package efp

import (
	"fmt"
	"testing"
)

func TestParseSimpleFieldValid(t *testing.T) {
	p, _ := PrototypeString("name : string")
	assertNow(t, len(p.fields) == 1, "wrong field length")
	// valid example
	e, errs := p.ValidateString(`name = "ender"`)
	assert(t, errs == nil, "errs should be nil")

	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, len(e.Field("name", 0).Values()) == 1, "wrong children number")
	assert(t, e.Field("name", 0).Value(0) == "ender", "xxx")

}

func TestParseSimpleFieldInvalid(t *testing.T) {
	p, _ := PrototypeString("name : string")

	// invalid example
	_, errs := p.ValidateString(`name = 6`)

	assert(t, errs != nil, "errs should not be nil")
}

func TestParseArrayFieldValid(t *testing.T) {
	p, _ := PrototypeString("name : [string]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, e.fields["name"] != nil, "fields should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "didn't find field")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")
	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")

	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")
}

func TestParseArrayFieldMinimumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, e.fields["name"] != nil, "fields should not be nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")
	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

}

func TestParseArrayFieldMinimumInvalid(t *testing.T) {
	p, _ := PrototypeString("name : [2:string]")
	// invalid
	_, errs := p.ValidateString(`name = ["ender"]`)
	assert(t, errs != nil, "errs should not be nil")
}

func TestParseArrayFieldMaximumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string:2]")
	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	//assert(t, p.fields["name"].Type(0, 1) == standards["string"].value, "wrong type input")
	assertNow(t, e.fields["name"] != nil, "fields should not be nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")
}

func TestParseArrayFieldMaximumInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string:2]")

	// invalid
	e, errs := p.ValidateString(`name = ["ender", "me", "him"]`)

	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")
}

func TestParseArrayFieldFixedValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string:2]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, e.fields["name"] != nil, "fields should not be nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")
}

func TestParseArrayFieldFixedInvalid(t *testing.T) {
	p, _ := PrototypeString("name : [2:string:2]")
	// invalid
	e, errs := p.ValidateString(`name = ["ender", "me", "him"]`)
	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = ["ender"]`)
	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")
}

func TestParseArrayFieldDisjunctionValid(t *testing.T) {

	// valid
	p, _ := PrototypeString("name : [string|int]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0 "+e.Field("name", 0).Value(0, 0))
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = [6, 7]`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")
	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionInvalid(t *testing.T) {

	p, _ := PrototypeString("name : [int|float]")

	// invalid
	e, errs := p.ValidateString(`name = [true, false]`)

	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldDisjunctionMinimumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string|int]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")

	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	//valid
	e, errs = p.ValidateString(`name = [6, 7]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionMinimumInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:float|int]")

	// invalid
	e, errs := p.ValidateString(`name = [true, false]`)
	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

	/* invalid
	//TODO: this halts??? e, errs = p.ValidateString(`name = [-0, -0]`)

	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")*/

	// invalid
	e, errs = p.ValidateString(`name = [6]`)
	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldDisjunctionMaximumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string|int:2]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	//valid
	e, errs = p.ValidateString(`name = [6, 7]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionMaximumInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string|int:2]")

	/* invalid
	e, errs := p.ValidateString(`name = [0.99, 0.22]`)

	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")*/

	// invalid
	e, errs := p.ValidateString(`name = ["a", "b", "c"]`)

	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = [6, 7, 8]`)

	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldDisjunctionFixedValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string|int:2]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	//valid
	e, errs = p.ValidateString(`name = [6, 7]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)

	assert(t, errs == nil, "errs should be nil")
	assertNow(t, e != nil, "e should not be nil")
	assertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	assertNow(t, e.Field("name", 0) != nil, "name nil")
	assertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	assertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionFixedInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:float|int:2]")

	// invalid
	e, errs := p.ValidateString(`name = [false, false]`)
	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = ["a", "b", "c"]`)
	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = [6, 7, 8]`)
	assert(t, errs != nil, "errs should not be nil")
	assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldTwoDimensionalDisjunction(t *testing.T) {
	// valid
	p, errs := PrototypeString("name : [2:[2:string|int:2]:2]")
	assert(t, errs == nil, "errs should be nil")
	e, errs := p.ValidateString(`name = [["ender", "me"], ["me", "ender"]]`)
	fmt.Println(errs)
	assert(t, len(e.Field("name", 0).Values(0, 0)) == 2, "wrong array length")
	assert(t, errs == nil, "errs should be nil")
	assert(t, e != nil, "e should not be nil")
}

func TestParseArrayFieldTwoDimensionalDisjunctionArrays(t *testing.T) {
	// valid
	p, errs := PrototypeString("name : [2:[2:string:2]|[2:int:2]:2]")
	assert(t, errs == nil, "errs should be nil")
	e, errs := p.ValidateString(`name = [["ender", "me"], ["me", "ender"]]`)
	fmt.Println(errs)
	assert(t, len(e.Field("name", 0).Values(0, 0)) == 2, "wrong array length")
	assert(t, errs == nil, "errs should be nil")
	assert(t, e != nil, "e should not be nil")
}
