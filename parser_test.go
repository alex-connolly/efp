package efp

import (
	"testing"

	"github.com/end-r/goutil"
)

func TestParseSimpleFieldValid(t *testing.T) {
	p, _ := PrototypeString("name : string")
	goutil.AssertNow(t, len(p.fields) == 1, "wrong field length")
	// valid example
	e, errs := p.ValidateString(`name = "ender"`)
	goutil.Assert(t, errs == nil, "errs should be nil")

	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, len(e.Field("name", 0).Values()) == 1, "wrong children number")
	goutil.Assert(t, e.Field("name", 0).Value(0) == "ender", "xxx")

}

func TestParseSimpleFieldInvalid(t *testing.T) {
	p, _ := PrototypeString("name : string")

	// invalid example
	_, errs := p.ValidateString(`name = 6`)

	goutil.Assert(t, errs != nil, "errs should not be nil")
}

func TestParseArrayFieldValid(t *testing.T) {
	p, _ := PrototypeString("name : [string]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, e.fields["name"] != nil, "fields should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "didn't find field")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")
	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")

	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")
}

func TestParseArrayFieldMinimumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, e.fields["name"] != nil, "fields should not be nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")
	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

}

func TestParseArrayFieldMinimumInvalid(t *testing.T) {
	p, _ := PrototypeString("name : [2:string]")
	// invalid
	_, errs := p.ValidateString(`name = ["ender"]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
}

func TestParseArrayFieldMaximumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string:2]")
	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	//goutil.Assert(t, p.fields["name"].Type(0, 1) == standards["string"].value, "wrong type input")
	goutil.AssertNow(t, e.fields["name"] != nil, "fields should not be nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")
}

func TestParseArrayFieldMaximumInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string:2]")

	// invalid
	e, errs := p.ValidateString(`name = ["ender", "me", "him"]`)

	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")
}

func TestParseArrayFieldFixedValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string:2]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, e.fields["name"] != nil, "fields should not be nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")
}

func TestParseArrayFieldFixedInvalid(t *testing.T) {
	p, _ := PrototypeString("name : [2:string:2]")
	// invalid
	e, errs := p.ValidateString(`name = ["ender", "me", "him"]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = ["ender"]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")
}

func TestParseArrayFieldDisjunctionValid(t *testing.T) {

	// valid
	p, _ := PrototypeString("name : [string|int]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0 "+e.Field("name", 0).Value(0, 0))
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = [6, 7]`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")
	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionInvalid(t *testing.T) {

	p, _ := PrototypeString("name : [int|float]")

	// invalid
	e, errs := p.ValidateString(`name = [true, false]`)

	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldDisjunctionMinimumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string|int]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")

	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	//valid
	e, errs = p.ValidateString(`name = [6, 7]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionMinimumInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:float|int]")

	// invalid
	e, errs := p.ValidateString(`name = [true, false]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

	/* invalid
	//TODO: this halts??? e, errs = p.ValidateString(`name = [-0, -0]`)

	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")*/

	// invalid
	e, errs = p.ValidateString(`name = [6]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldDisjunctionMaximumValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string|int:2]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	//valid
	e, errs = p.ValidateString(`name = [6, 7]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionMaximumInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [string|int:2]")

	/* invalid
	e, errs := p.ValidateString(`name = [0.99, 0.22]`)

	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")*/

	// invalid
	e, errs := p.ValidateString(`name = ["a", "b", "c"]`)

	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = [6, 7, 8]`)

	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldDisjunctionFixedValid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:string|int:2]")

	e, errs := p.ValidateString(`name = ["ender", "me"]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "me", "invalid value 1")

	//valid
	e, errs = p.ValidateString(`name = [6, 7]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "6", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "7", "invalid value 1")

	// valid
	e, errs = p.ValidateString(`name = ["ender", 6]`)

	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, len(e.fields["name"]) == 1, "field length wrong")
	goutil.AssertNow(t, e.Field("name", 0) != nil, "name nil")
	goutil.AssertNow(t, e.Field("name", 0).Values() != nil, "values should not be nil")

	goutil.AssertNow(t, len(e.Field("name", 0).Values(0)) == 2, "wrong value length")
	goutil.Assert(t, e.Field("name", 0).Value(0, 0) == "ender", "invalid value 0")
	goutil.Assert(t, e.Field("name", 0).Value(0, 1) == "6", "invalid value 1")
}

func TestParseArrayFieldDisjunctionFixedInvalid(t *testing.T) {
	// valid
	p, _ := PrototypeString("name : [2:float|int:2]")

	// invalid
	e, errs := p.ValidateString(`name = [false, false]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = ["a", "b", "c"]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

	// invalid
	e, errs = p.ValidateString(`name = [6, 7, 8]`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	goutil.Assert(t, e != nil, "e should not be nil")

}

func TestParseArrayFieldTwoDimensionalDisjunction(t *testing.T) {
	// valid
	p, errs := PrototypeString("name : [2:[2:string|int:2]:2]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	e, errs := p.ValidateString(`name = [["ender", "me"], ["me", "ender"]]`)
	goutil.Assert(t, len(e.Field("name", 0).Values(0)) == 2, "wrong length 0")
	goutil.Assert(t, len(e.Field("name", 0).Values(0, 0)) == 2, "wrong length 0 0")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.Assert(t, e != nil, "e should not be nil")
}

func TestParseArrayFieldTwoDimensionalDisjunctionArrays(t *testing.T) {
	// valid
	p, errs := PrototypeString("name : [2:[2:string:2]|[2:int:2]:2]")
	goutil.Assert(t, errs == nil, "errs should be nil")
	e, errs := p.ValidateString(`name = [["ender", "me"], ["me", "ender"]]`)
	goutil.Assert(t, len(e.Field("name", 0).Values(0, 0)) == 2, "wrong array length")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.Assert(t, e != nil, "e should not be nil")
}

func TestParseFieldAlias(t *testing.T) {
	p, errs := PrototypeString(`alias x = name : string   x  `)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.Assert(t, p.Field("name").TypeValue(0) == standards["string"], "wrong value")
}

func TestParseElement(t *testing.T) {
	p, errs := PrototypeString(`name {
			name : string
		}`)
	goutil.AssertNow(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, len(p.Element("name").fields) == 1, "wrong pfield length")
	goutil.AssertNow(t, p.Element("name").Field("name") != nil, "field is nil")
	e, errs := p.ValidateString(`name {
		name = "ender"
	}`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, e.Element("name", 0).Field("name", 0) != nil, "shouldn't be nil")
	goutil.AssertNow(t, e.Element("name", 0).Field("name", 0).Value(0) == "ender", "wrong value")
}

func TestParseElementParameters(t *testing.T) {
	p, errs := PrototypeString(`name(string, string, int) {
			name : string
		}`)
	e, errs := p.ValidateString(`name("hi", "i'm", 16){ name = "ender" }`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.AssertNow(t, e.Element("name", 0) != nil, "element is nil")
	goutil.AssertNow(t, e.Element("name", 0).Parameters() != nil, "parameters are nil")
	goutil.AssertNow(t, e.Element("name", 0).Parameter(2).Value() == "16", "wrong parameter value")
	goutil.AssertNow(t, e.Element("name", 0).Field("name", 0) != nil, "Field is nil")
	goutil.AssertNow(t, e.Element("name", 0).Field("name", 0).Value(0) == "ender", "wrong value")
}
