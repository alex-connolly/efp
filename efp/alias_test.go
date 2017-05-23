package efp

import (
	"regexp"
	"testing"
)

func TestAliasingStandardAliases(t *testing.T) {
	p := createPrototypeParserString("")
	assert(t, len(p.prototype.textAliases) == len(standards), "wrong number of standard aliases")
}

func TestAliasingStandardInt(t *testing.T) {
	p := createPrototypeParserString("")
	assertNow(t, p.prototype.textAliases != nil, "text aliases is nil")
	regex, err := regexp.Compile(p.prototype.textAliases["int"].value)
	assertNow(t, err == nil, "error compiling int regex")
	assert(t, regex.MatchString("99"), "int regex didn't match")
	assert(t, regex.MatchString("0"), "int regex didn't match")
	assert(t, regex.MatchString("-100"), "int regex didn't match")

	assert(t, !regex.MatchString("aaa"), "int regex did match")
	assert(t, !regex.MatchString("-0"), "int regex did match")
}

func TestAliasingStandardUInt(t *testing.T) {
	p := createPrototypeParserString("")
	regex, err := regexp.Compile(p.prototype.textAliases["uint"].value)
	assertNow(t, err == nil, "error compiling uint regex")
	assert(t, regex.MatchString("99"), "uint regex didn't match")
	assert(t, regex.MatchString("0"), "uint regex didn't match")

	assert(t, !regex.MatchString("-100"), "uint regex did match")
	assert(t, !regex.MatchString("-0"), "uint regex did match")
}

func TestAliasingStandardFloat(t *testing.T) {
	p := createPrototypeParserString("")
	regex, err := regexp.Compile(p.prototype.textAliases["float"].value)
	assertNow(t, err == nil, "error compiling float regex")
	assert(t, regex.MatchString("99"), "float regex didn't match")
	assert(t, regex.MatchString("0"), "float regex didn't match")
	assert(t, regex.MatchString("-100"), "float regex didn't match")
	assert(t, regex.MatchString("0.5"), "float regex didn't match")

	assert(t, !regex.MatchString("aaa"), "float regex did match")
	assert(t, !regex.MatchString("-0"), "float regex did match")
}

func TestAliasingStandardBool(t *testing.T) {
	p := createPrototypeParserString("")
	regex, err := regexp.Compile(p.prototype.textAliases["bool"].value)
	assertNow(t, err == nil, "error compiling bool regex")
	assert(t, regex.MatchString("true"), "bool regex didn't match")
	assert(t, regex.MatchString("false"), "bool regex didn't match")

	assert(t, !regex.MatchString("tru"), "bool regex did match")
	assert(t, !regex.MatchString("flase"), "bool regex did match")
}

func TestAliasingFieldAlias(t *testing.T) {
	p, errs := PrototypeString(`
        alias x = name : string
        x`)
	assertNow(t, errs == nil, "errs are not nil")
	assertNow(t, p.fields["name"] != nil, "name is nil")
}

func TestAliasingElementAlias(t *testing.T) {
	p, errs := PrototypeString(`
        alias x = name {

		}
        x`)
	assertNow(t, errs == nil, "errs are not nil")
	assertNow(t, p.elements != nil, "elements is nil")
	assertNow(t, p.elements["name"] != nil, "name is nil")
}

func TestAliasingTextAliasMax(t *testing.T) {
	const limit = 2
	const regex = "[a-z]{3}"
	p, errs := PrototypeString(`
        alias LIMIT = 2
		<LIMIT:"[a-z]{3}":LIMIT> : [LIMIT:string:LIMIT]
		`)
	assertNow(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields[regex] != nil, "field is nil")
	assertNow(t, p.fields[regex].types[0].isArray, "not array")
	assertNow(t, p.fields[regex].types[0].types[0].value.String() == standards["string"].value, "incorrect regex")
	assertNow(t, p.fields[regex].types[0].max == limit, "incorrect value max")
	assertNow(t, p.fields[regex].types[0].min == limit, "incorrect value min")
	assertNow(t, p.fields[regex].key.min == limit, "incorrect key min")
	assertNow(t, p.fields[regex].key.max == limit, "incorrect key max")
}

func TestAliasingTextAliasValue(t *testing.T) {
	p, errs := PrototypeString(`
        alias x = string
        name : x`)
	assertNow(t, errs == nil, "errs should be nil")
	assertNow(t, p.fields["name"] != nil, "name is nil")
	assertNow(t, p.fields["name"].types[0].value.String() == standards["string"].value, "wrong value")
}

func TestAliasingDoubleIndirection(t *testing.T) {
	p := createPrototypeParserString(`
        alias y = string
        alias x = name : y
        x`)
	assertNow(t, p.prototype.fields["name"] != nil, "name is nil")
	assertNow(t, p.prototype.fields["name"].types[0].value.String() == standards["string"].value, "wrong value")
}

// test that element recursion is allowed
func TestAliasingRecursionValid(t *testing.T) {
	p, errs := PrototypeString(`
        alias p = x {
            p
        }
        p`)
	assertNow(t, errs == nil, "errs should be nil")
	assertNow(t, p.elements["x"] != nil, "x is nil")
}

// test that field recursion is disallowed
func TestAliasingRecursionInvalid(t *testing.T) {
	_, errs := PrototypeString(`alias x = x
		x`)
	assertNow(t, errs != nil, "errs should not be nil")
}
