package efp

type tokenType int

const (
	tknValue = iota
	tknNumber
	tknComment
	tknAlias
	tknAssign
	tknComma
	tknOpenBrace
	tknCloseBrace
	tknOpenSquare
	tknCloseSquare
	tknOpenBracket
	tknCloseBracket
	tknOpenCorner
	tknCloseCorner
	tknRequired
	tknColon
	tknOr
	tknNone
)

func getProtoTokens() []protoToken {
	return []protoToken{
		protoToken{"Open Square Bracket", is('['), processOperator(tknOpenSquare)},
		protoToken{"Close Square Bracket", is(']'), processOperator(tknCloseSquare)},
		protoToken{"Open Bracket", is('('), processOperator(tknOpenBracket)},
		protoToken{"Close Bracket", is(')'), processOperator(tknCloseBracket)},
		protoToken{"Open Brace", is('{'), processOperator(tknOpenBrace)},
		protoToken{"Close Brace", is('}'), processOperator(tknCloseBrace)},
		protoToken{"Open Corner", is('<'), processOperator(tknOpenCorner)},
		protoToken{"Close Corner", is('>'), processOperator(tknCloseCorner)},
		protoToken{"Assignment Operator", is('='), processOperator(tknAssign)},
		protoToken{"Required Operator", is('!'), processOperator(tknRequired)},
		protoToken{"Comma", is(','), processOperator(tknComma)},
		protoToken{"Colon", is(':'), processOperator(tknColon)},
		protoToken{"Or", is('|'), processOperator(tknOr)},
		protoToken{"New Line", isNewLine, processNewLine},
		protoToken{"Whitespace", isWhitespace, processIgnored},
		protoToken{"String", isString, processString},
		protoToken{"Number", isNumber, processNumber},
		protoToken{"Identifier", isIdentifier, processIdentifier},
	}
}

func processOperator(tkn tokenType) processorFunc {
	return func(l *lexer) (t token) {

		t.start = l.offset
		t.end = l.offset
		t.tkntype = tkn
		l.offset++
		return t
	}
}

func isIdentifier(b byte) bool {
	return ('A' <= b && b <= 'z') || ('0' <= b && b <= '9') || (b == '_')
}

func isNumber(b byte) bool {
	return ('0' <= b && b <= '9')
}

func isString(b byte) bool {
	return ((b == '"') || (b == '\''))
}

func isWhitespace(b byte) bool {
	return (b == ' ') || (b == '\t') || (b == 'r')
}

func isNewLine(b byte) bool {
	return (b == '\n')
}

func is(a byte) isFunc {
	return func(b byte) bool {
		return b == a
	}
}

type isFunc func(byte) bool
type processorFunc func(*lexer) token

type protoToken struct {
	name       string // for debugging
	identifier isFunc
	process    processorFunc
}

type token struct {
	tkntype tokenType
	start   int
	end     int
}

func checkAndReplaceKeywordRegex(value string) string {
	switch value {
	case "id":
		return "[a-zA-Z_]+"
	case "string":
		return "\"[^()]\""
	case "int":
		return "[0-9]+"
	case "bool":
		return "true|false"
	case "float":
		return "[0-9]*.[0-9]+"
	}
	return value
}
