package efp

// returns the distance between two tokens, but with:
// string|int|[string:3]|[[2:int]] == 1
// alias x = 2
// string|int|x = 1
// string | int | x = 1
// ALIAS ALIAS2 = 2
// first token is 0 away, second is 1 away etc...
// horrific implementation currently
func realDistance(p *parser, tk tokenType, number int) int {
	found := 0
	count := 0
	inValue := false
	var prev tokenType
	prev = tknNone
	for i, t := range p.lexer.tokens {
		if !inValue {
			if t.tkntype == tk {
				found++
				if found == number {
					return count
				}
			}
			switch t.tkntype {
			case tknValue, tknString, tknNumber:
				//fmt.Printf("in value\n")
				inValue = true
			}
			count++
		} else {
			switch t.tkntype {
			case tknValue, tknString, tknNumber:
				switch prev {
				case tknValue, tknString, tknNumber:
					if t.tkntype == tk {
						found++
						if found == number {
							return count
						}
					}
					count++
					break
				}
				break
			case tknColon:
				// [2:string] --> ignore
				// [string:2] --> ignore
				// x : string --> keep
				if i < len(p.lexer.tokens)-1 && i > 0 {
					if p.lexer.tokens[i+1].tkntype != tknNumber && p.lexer.tokens[i-1].tkntype != tknNumber {
						//fmt.Printf("%s --> hi colon xxx %d\n", p.lexer.buffer, i)
						inValue = false
						if t.tkntype == tk {
							found++
							if found == number {
								return count
							}
						}
						count++
					} else {
						//fmt.Printf("value colon\n ")
					}
				}
			case tknOpenSquare, tknCloseSquare, tknOr:
				// do nothing
				break
			default:
				//fmt.Printf("value ended with: %d\n", t.tkntype)
				inValue = false
				if t.tkntype == tk {
					found++
					if found == number {
						return count
					}
				}
				count++
			}
		}
		prev = t.tkntype
	}
	return -1
}

// field of the form key = value
func isField(p *parser) bool {
	// field can be in one of these forms:
	// key = value
	return (realDistance(p, tknValue, 1) == 0 && realDistance(p, tknAssign, 1) == 1) ||
		// "key" = value
		(realDistance(p, tknString, 1) == 0 && realDistance(p, tknAssign, 1) == 1)

}

// elements are of the form key { or key(params){
func isElement(p *parser) bool {
	// key {}
	return (realDistance(p, tknValue, 1) == 0 && realDistance(p, tknOpenBrace, 1) == 1) ||
		// key(params){
		(realDistance(p, tknValue, 1) == 0 && realDistance(p, tknOpenBracket, 1) == 1) ||
		//"key"{}
		(realDistance(p, tknString, 1) == 0 && realDistance(p, tknOpenBrace, 1) == 1) ||
		// "key"(params){
		(realDistance(p, tknString, 1) == 0 && realDistance(p, tknOpenBracket, 1) == 1)
}

// closures are just }
func isElementClosure(p *parser) bool {
	return p.lexer.tokens[p.index].tkntype == tknCloseBrace
}

// must be run last to exclude other possibilities
func isTextAlias(p *parser) bool {
	// any other alias will fit in this category
	return isAlias(p) && realDistance(p, tknAssign, 1) == 2
}

// alias x : key = value
func isFieldAlias(p *parser) bool {
	return isAlias(p) && isPrototypeFieldWithOffset(p, 3, 2)
}

func isAlias(p *parser) bool {
	return p.lexer.tokenString(p.peek(p.index)) == "alias" &&
		realDistance(p, tknValue, 2) == 1 &&
		realDistance(p, tknAssign, 1) == 2
}

// alias divs = divs(){}
func isElementAlias(p *parser) bool {
	return isAlias(p) && isPrototypeElementWithOffset(p, 3, 2)
}

func isPrototypeField(p *parser) bool {
	return isPrototypeFieldWithOffset(p, 0, 0)
}

// extra is the number of "extra" values (alias x =) = 2
func isPrototypeFieldWithOffset(p *parser, offset int, extra int) bool {
	// key : value
	return (realDistance(p, tknValue, 1+extra) == offset && realDistance(p, tknColon, 1) == 1+offset) ||
		// "key" : value
		(realDistance(p, tknString, 1) == offset && realDistance(p, tknColon, 1) == 1+offset) ||
		// <key> : values || <3:key> : value || <key:3> : value || <3:key:3> : value
		(realDistance(p, tknOpenCorner, 1) == offset && realDistance(p, tknColon, 1) == 3+offset)
}

func isPrototypeElement(p *parser) bool {
	return isPrototypeElementWithOffset(p, 0, 0)
}

// extra is the number of "extra" values (alias x =) = 2
func isPrototypeElementWithOffset(p *parser, offset int, extra int) bool {
	// key {}
	return (realDistance(p, tknValue, 1+extra) == offset && realDistance(p, tknOpenBrace, 1) == 1+offset) ||
		// <key>{}
		(realDistance(p, tknOpenCorner, 1) == offset && realDistance(p, tknOpenBrace, 1) == 3+offset) ||
		// <3:string|int>{}
		(realDistance(p, tknOpenCorner, 1) == offset && realDistance(p, tknOpenBrace, 1) == 5+offset) ||
		// <3:string|int:3>{}
		(realDistance(p, tknOpenCorner, 1) == offset && realDistance(p, tknOpenBrace, 1) == 7+offset) ||
		// key(){}
		(realDistance(p, tknValue, 1+extra) == offset && realDistance(p, tknOpenBracket, 1) == 1+offset) ||
		// <key>(){}
		(realDistance(p, tknOpenCorner, 1) == offset && realDistance(p, tknOpenBracket, 1) == 3+offset) ||
		// <3:string>(){} or <string:3>(){}
		(realDistance(p, tknOpenCorner, 1) == offset && realDistance(p, tknOpenBracket, 1) == 5+offset) ||
		// <3:string|int:3>(){}, <3:string|"[A-Z]+"|name:3>(){}
		(realDistance(p, tknOpenCorner, 1) == offset && realDistance(p, tknOpenBracket, 1) == 7+offset)

}

func isDiscoveredAlias(p *parser) bool {
	return realDistance(p, tknValue, 1) == 0
}
