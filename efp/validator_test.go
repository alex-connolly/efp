package efp

import "testing"

func TestValidateField(t *testing.T) {

}

func TestValidateKey(t *testing.T) {
	p := basicParser("name : string")
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields["name"] != nil, "fields is nil")
	// valid example
	key := p.validateKey("name")
	assertNow(t, p.errs == nil, "Errors not nil")
	assert(t, key == "name", "Failed equality test")

	p = basicParser(`"[a-z]+" : string`)
	parsePrototypeField(p)
	assertNow(t, p.prototype.fields["[a-z]+"] != nil, "fields is nil")
	// valid example
	key = p.validateKey("name")
	assertNow(t, p.errs == nil, "Errors not nil")
	assert(t, key == "[a-z]+", "Failed equality test")
}
