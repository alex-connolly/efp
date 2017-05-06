package efp

type lexer struct {
	buffer    []byte
	offset    int
	line      int
	column    int
	tokens    []token
	numTokens int
	length    int
}

func (l *lexer) next() {
	if l.offset == l.length {
		return
	}
	b := l.nextByte()
	for _, pt := range getProtoTokens() {
		if pt.identifier(b) {
			t := pt.process(l)
			if t.tkntype != tknNone {
				l.tokens = append(l.tokens, t)
			}
			break
		}
	}
	l.next()
}

func (l *lexer) isEOF() bool {
	return l.offset == l.length
}

func (l *lexer) tokenString(t token) string {
	data := make([]byte, t.end-t.start)
	copy(data, l.buffer[t.start-1:t.end-1])
	return string(data)
}

func (l *lexer) nextByte() byte {
	b := l.buffer[l.offset]
	l.offset++
	return b
}

func lex(bytes []byte) *lexer {
	l := new(lexer)
	l.buffer = bytes
	l.length = len(bytes)
	l.next()
	return l
}
