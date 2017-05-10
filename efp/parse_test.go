package efp

import "testing"

func TestParseSimpleField(t *testing.T) {
	p := basicParser("name : string")
	// valid example
	ast, errs := p.parseBytes([]byte(`name = "ender"`))
	assert(t, errs == nil, "errs should be nil")
	assert(t, ast != nil, "ast should not be nil")
	assert(t, ast.fields["name"].value == "ender", "invalid value")
	// invalid example
	ast, errs = p.parseBytes([]byte(`name = ender`))
	assert(t, errs != nil, "errs should not be nil")
}

func TestParseArrayField(t *testing.T) {
	p := basicParser("name : [string]")
	ast, errs := p.parseBytes([]byte(`name = ["ender", "me"]`))
	assert(t, errs == nil, "errs should be nil")
	assert(t, ast != nil, "ast should not be nil")
	assert(t, ast.fields["name"].value.children[0] == "ender", "invalid value 0")
	assert(t, ast.fields["name"].value.children[1] == "me", "invalid value 1")
}

func TestParseArrayFieldMinimum(t *testing.T) {

}
