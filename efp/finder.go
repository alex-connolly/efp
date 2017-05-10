package efp

import "fmt"

const alias = "alias"

// returns the distance between two tokens, but with:
// string|int|[string]|[[int]] == 1
func realDistance(p *parser, tk tokenType) int {
	count := 0
	inValue := false
	for _, t := range p.lexer.tokens {
		if t.tkntype == tk {
			return count
		}
		if t.tkntype == tknValue ||
			t.tkntype == tknOr ||
			t.tkntype == tknOpenSquare ||
			t.tkntype == tknCloseSquare {
			if !inValue {
				count++
				inValue = true
			}
		} else {
			inValue = false
			count++
		}
	}
	return -1
}

// field of the form key = value
func isField(p *parser) bool {
	// field can be in one of these forms:
	// key = value
	return (realDistance(p, tknValue) == 0 && realDistance(p, tknAssign) == 1)

}

// elements are of the form key { or key(params){
func isElement(p *parser) bool {
	// key {}
	return (realDistance(p, tknValue) == 0 && realDistance(p, tknOpenBrace) == 1) &&
		// key(params){
		(realDistance(p, tknValue) == 0 && realDistance(p, tknOpenBracket) == 1)
}

// closures are just }
func isElementClosure(p *parser) bool {
	return realDistance(p, tknCloseBrace) == 0
}

// must be run last to exclude other possibilities
func isTextAlias(p *parser) bool {
	return realDistance(p, tknValue) == 0
}

// alias x : key = value
func isFieldAlias(p *parser) bool {
	return p.lexer.tokenString(p.peek(p.index)) == alias &&
		p.peek(p.index+2).tkntype == tknColon &&
		p.peek(p.index+4).tkntype == tknAssign
}

// alias divs : divs(){}
func isElementAlias(p *parser) bool {
	return (p.lexer.tokenString(p.peek(p.index)) == alias &&
		p.peek(p.index+1).tkntype == tknValue &&
		p.peek(p.index+2).tkntype == tknColon)
}

func isPrototypeField(p *parser) bool {
	fmt.Printf("d: %d, %d\n", realDistance(p, tknValue), realDistance(p, tknColon))
	// key : value
	return (realDistance(p, tknValue) == 0 && realDistance(p, tknColon) == 1) ||
		// <key> : value
		(realDistance(p, tknOpenCorner) == 0 && realDistance(p, tknColon) == 3) ||
		// <3:key> : value
		(realDistance(p, tknOpenCorner) == 0 && realDistance(p, tknColon) == 5) ||
		// <key:3> : value is the same as above
		// <3:key:5> : value
		(realDistance(p, tknOpenCorner) == 0 && realDistance(p, tknColon) == 7)
}

// currently won't work:
// <3:int|string:3>(){}
func isPrototypeElement(p *parser) bool {

	// key {}
	return (realDistance(p, tknValue) == 0 && realDistance(p, tknOpenBrace) == 1) ||
		//key(){}
		(realDistance(p, tknValue) == 0 && realDistance(p, tknOpenBracket) == 1) ||
		// <key>{}
		(realDistance(p, tknOpenCorner) == 0 && realDistance(p, tknOpenBrace) == 3) ||
		// <key|k>(){}
		(realDistance(p, tknOpenCorner) == 0 && realDistance(p, tknOpenBracket) == 3) ||
		// <3:string|int>{}
		(realDistance(p, tknOpenCorner) == 0 && realDistance(p, tknOpenBrace) == 6)
}

func isAlias(p *parser) bool {
	return realDistance(p, tknValue) == 0
}
