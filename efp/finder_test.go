package efp

import "testing"

func TestIsPrototypeField(t *testing.T) {
	p := &parser{lexer: lex([]byte("name : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("name : string!"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("<name> : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("<3:name> : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("<name:3> : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("<3:name:3> : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("<name|string> : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("<3:name|string> : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte("<name|string:3> : string"))}
	failIf(t, !isPrototypeField(p))
	p = &parser{lexer: lex([]byte(`<3:name|string|"a-zA-z{5}":3> : string`))}
	failIf(t, !isPrototypeField(p))
}

func TestIsPrototypeElement(t *testing.T) {
	p := &parser{lexer: lex([]byte("name {}"))}
	failIf(t, !isPrototypeElement(p))
	p = &parser{lexer: lex([]byte("<name> {}"))}
	failIf(t, !isPrototypeElement(p))
	p = &parser{lexer: lex([]byte("<3:name> {}"))}
	failIf(t, !isPrototypeElement(p))
	p = &parser{lexer: lex([]byte("<name:3> {}"))}
	failIf(t, !isPrototypeElement(p))
	p = &parser{lexer: lex([]byte("<3:name:3> {}"))}
	failIf(t, !isPrototypeElement(p))
}

func TestIsFieldAlias(t *testing.T) {
	p := &parser{lexer: lex([]byte("alias x = name : string"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte("alias x = name : string!"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name> : string"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name> : string"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name:3> : string"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name:3> : string"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name|int> : string"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name|int|string:3> : string"))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte(`alias x = <3:"A-Za-z":3> : string`))}
	failIf(t, !isFieldAlias(p))
	p = &parser{lexer: lex([]byte(`alias x = <3:"A-Za-z"|int:3> : string`))}
	failIf(t, !isFieldAlias(p))
}

func TestIsElementAlias(t *testing.T) {
	p := &parser{lexer: lex([]byte("alias x = name {}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name> {}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name> {}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name:3> {}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name:3> {}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = name(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name:3>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name:3>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name|int|string>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name|int|string>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <name|int|string:3>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name|int|string:3>(int){}"))}
	failIf(t, !isElementAlias(p))
	p = &parser{lexer: lex([]byte("alias x = <3:name|int|string:3>(int){}"))}
	failIf(t, !isElementAlias(p))
}

func TestIsField(t *testing.T) {
	p := &parser{lexer: lex([]byte("name = 6"))}
	failIf(t, !isField(p))
	p = &parser{lexer: lex([]byte(`name = "www"`))}
	failIf(t, !isField(p))
	p = &parser{lexer: lex([]byte("name = hi"))}
	failIf(t, !isField(p))
	p = &parser{lexer: lex([]byte(`name = [hi, me, c]`))}
	failIf(t, !isField(p))
	p = &parser{lexer: lex([]byte(`name = ["hi", "me", "c"]`))}
	failIf(t, !isField(p))
	p = &parser{lexer: lex([]byte(`name = [["hi", "me"], [6, 5, 7]]`))}
	failIf(t, !isField(p))
}

func TestIsElement(t *testing.T) {
	p := &parser{lexer: lex([]byte("name {}"))}
	failIf(t, !isElement(p))
	p = &parser{lexer: lex([]byte("name(int){}"))}
	failIf(t, !isElement(p))
	p = &parser{lexer: lex([]byte("name(){}"))}
	failIf(t, !isElement(p))
}

func TestIsOperator(t *testing.T) {
	if !is(',')(',') {
		t.Fail()
	}
}
