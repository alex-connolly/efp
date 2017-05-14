package efp

import "testing"

func TestParseSimpleField(t *testing.T) {
	p := basicParser("name : string")
	// valid example
	p.parseString(`name = "ender"`)
	parseField(p)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assertNow(t, len(p.scope.fields["name"][0].value.children) == 1, "wrong children number")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value "+p.scope.fields["name"][0].value.children[0].value)
	// invalid example
	p.parseString(`name = ender`)
	assert(t, p.errs != nil, "errs should not be nil")
}

func TestParseArrayField(t *testing.T) {
	p := new(parser)
	p.prototypeString("name : [string]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")
}

func TestParseArrayFieldMinimum(t *testing.T) {
	// valid
	p := basicParser("name : [2:string]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	// invalid
	p.parseString(`name = ["ender", "me", "him"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldMaximum(t *testing.T) {
	// valid
	p := basicParser("name : [string:2]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	// invalid
	p.parseString(`name = ["ender", "me", "him"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldFixed(t *testing.T) {
	// valid
	p := basicParser("name : [2:string:2]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, p.scope.fields["name"] != nil, "fields should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	// invalid
	p.parseString(`name = ["ender", "me", "him"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseString(`name = ["ender"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldDisjunction(t *testing.T) {
	// valid
	p := basicParser("name : [string|int]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseString(`name = [6, 7]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseString(`name = ["ender", 6]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseString(`name = [hello, 6]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionMinimum(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseString(`name = [6, 7]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseString(`name = ["ender", 6]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseString(`name = [hello, 6]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseString(`name = ["a"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseString(`name = [6]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionMaximum(t *testing.T) {
	// valid
	p := basicParser("name : [string|int:2]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseString(`name = [6, 7]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseString(`name = ["ender", 6]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseString(`name = [hello, 6]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseString(`name = ["a", "b", "c"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseString(`name = [6, 7, 8]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionFixed(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int:2]")
	p.parseString(`name = ["ender", "me"]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseString(`name = [6, 7]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseString(`name = ["ender", 6]`)
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assertNow(t, len(p.scope.fields["name"]) == 1, "field length wrong")
	assertNow(t, p.scope.fields["name"][0] != nil, "name nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseString(`name = [hello, 6]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseString(`name = ["a", "b", "c"]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseString(`name = [6, 7, 8]`)
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldTwoDimensionalDisjunction(t *testing.T) {
	// valid
	p := basicParser("name : [2:[2:string|int:2]:2]")
	p.parseString(`name = [["ender", "me"], ["me", "ender"]]`)
	assert(t, p.errs == nil, "errs should be nil")
}

func TestParseArrayFieldTwoDimensionalDisjunctionArrays(t *testing.T) {
	// valid
	p := basicParser("name : [2:[2:string:2|[2:int:2]:2]")
	p.parseString(`name = [["ender", "me"], ["me", "ender"]]`)
	assert(t, p.errs == nil, "errs should be nil")
}
