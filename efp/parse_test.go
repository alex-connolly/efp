package efp

import "testing"

func TestParseField(t *testing.T) {
	p := basicParser("name : string")
	ast, errs := p.parseBytes([]byte(`name = "xxx"`))
	if errs != nil {
		t.Fail()
	}
	if ast.fields["name"].value != "xxx" {
		t.Fail()
	}
}
