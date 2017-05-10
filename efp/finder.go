package efp

const alias = "alias"

// field of the form key = value
func isField(p *parser) bool {
	// field can be in one of these forms:
	// key = value
	return (p.peek(p.index+1).tkntype == tknAssign)

}

// elements are of the form key { or key(params){
func isElement(p *parser) bool {
	return (p.peek(p.index).tkntype == tknValue &&
		p.peek(p.index+1).tkntype == tknOpenBrace) ||
		(p.peek(p.index).tkntype == tknValue &&
			p.peek(p.index+1).tkntype == tknOpenBracket)
}

// closures are just }
func isElementClosure(p *parser) bool {
	return p.peek(p.index).tkntype == tknCloseBrace
}

// must be run last to exclude other possibilities
func isTextAlias(p *parser) bool {
	return p.peek(p.index).tkntype == tknValue
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
	// key : value
	return (p.peek(p.index).tkntype == tknValue && p.peek(p.index+1).tkntype == tknColon) ||
		// <key> : value
		(p.peek(p.index).tkntype == tknOpenCorner && p.peek(p.index+3).tkntype == tknColon) ||
		// <3:key> : value
		(p.peek(p.index).tkntype == tknOpenCorner && p.peek(p.index+5).tkntype == tknColon) ||
		// <key:3> : value is the same as above
		// <3:key:5> : value
		(p.peek(p.index).tkntype == tknOpenCorner && p.peek(p.index+7).tkntype == tknColon)
}

func isPrototypeElement(p *parser) bool {
	// key {}
	return (p.peek(p.index) == tknValue && p.peek(p.index+1).tkntype == tknOpenBrace) ||
		//key(){}
		(p.peek(p.index) == tknValue && p.peek(p.index+1).tkntype == tknOpenBracket) ||
        // <key>{}
        (p.peek(p.index) == tknValue && p.peek(p.index+3).tkntype == tknOpenBrace) ||
}
