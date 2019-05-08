package lex

// attempts to read a string token (literal contained in double quotes)
// expects the beginning double quote to have been consumed already
// consumes all bytes up to and including the terminating double quote
// TODO: DOES NOT CURRENTLY SUPPORT ESCAPE SEQUENCES
func (l *Lexer) readStringBytes() bool {
	for {
		b, err := l.r.ReadByte()
		if err != nil {
			return false
		}
		if b == '"' {
			return true
		}
		l.lastTokenBytes = append(l.lastTokenBytes, b)
	}
}
