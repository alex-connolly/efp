package efp

import "testing"

func TestParseSimpleField(t *testing.T) {
	p := basicParser("name : string")
	// valid example
	p.parseBytes([]byte(`name = "ender"`))
	assert(t, p.errs == nil, "errs should be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value")
	// invalid example
	p.parseBytes([]byte(`name = ender`))
	assert(t, p.errs != nil, "errs should not be nil")
}

func TestParseArrayField(t *testing.T) {
	p := basicParser("name : [string]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")
}

func TestParseArrayFieldMinimum(t *testing.T) {
	// valid
	p := basicParser("name : [2:string]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	// invalid
	p.parseBytes([]byte(`name = ["ender", "me", "him"]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldMaximum(t *testing.T) {
	// valid
	p := basicParser("name : [string:2]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	// invalid
	p.parseBytes([]byte(`name = ["ender", "me", "him"]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldFixed(t *testing.T) {
	// valid
	p := basicParser("name : [2:string:2]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	// invalid
	p.parseBytes([]byte(`name = ["ender", "me", "him"]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseBytes([]byte(`name = ["ender"]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")
}

func TestParseArrayFieldDisjunction(t *testing.T) {
	// valid
	p := basicParser("name : [string|int]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseBytes([]byte(`name = [6, 7]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseBytes([]byte(`name = ["ender", 6]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseBytes([]byte(`name = [hello, 6]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionMinimum(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseBytes([]byte(`name = [6, 7]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseBytes([]byte(`name = ["ender", 6]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseBytes([]byte(`name = [hello, 6]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseBytes([]byte(`name = ["a"]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseBytes([]byte(`name = [6]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionMaximum(t *testing.T) {
	// valid
	p := basicParser("name : [string|int:2]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseBytes([]byte(`name = [6, 7]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseBytes([]byte(`name = ["ender", 6]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseBytes([]byte(`name = [hello, 6]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseBytes([]byte(`name = ["a", "b", "c"]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseBytes([]byte(`name = [6, 7, 8]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldDisjunctionFixed(t *testing.T) {
	// valid
	p := basicParser("name : [2:string|int:2]")
	p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "me", "invalid value 1")

	//valid
	p.parseBytes([]byte(`name = [6, 7]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "6", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "7", "invalid value 1")

	// valid
	p.parseBytes([]byte(`name = ["ender", 6]`))
	assert(t, p.errs == nil, "errs should be nil")
	assertNow(t, p.scope != nil, "p.scope should not be nil")
	assert(t, p.scope.fields["name"][0].value.children[0].value == "ender", "invalid value 0")
	assert(t, p.scope.fields["name"][0].value.children[1].value == "6", "invalid value 1")

	// invalid
	p.parseBytes([]byte(`name = [hello, 6]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseBytes([]byte(`name = ["a", "b", "c"]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

	// invalid
	p.parseBytes([]byte(`name = [6, 7, 8]`))
	assert(t, p.errs != nil, "errs should not be nil")
	assert(t, p.scope != nil, "p.scope should not be nil")

}

func TestParseArrayFieldTwoDimensionalDisjunction(t *testing.T) {
	// valid
	p := basicParser("name : [2:[2:string|int:2]:2]")
	p.parseBytes([]byte(`name = [["ender", "me"], ["me", "ender"]]`))
	assert(t, p.errs == nil, "errs should be nil")
}

func TestParseArrayFieldTwoDimensionalDisjunctionArrays(t *testing.T) {
	// valid
	p := basicParser("name : [2:[2:string:2|[2:int:2]:2]")
	p.parseBytes([]byte(`name = [["ender", "me"], ["me", "ender"]]`))
	assert(t, p.errs == nil, "errs should be nil")
}
