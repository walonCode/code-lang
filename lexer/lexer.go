package lexer

import "github.com/walonCode/code-lang/token"

type Lexer struct {
	input string
	position int
	readPosition int
	ch byte
}

//methods on the lexer
func (l *Lexer)readChar(){
	//check if we reach the end of the input
	if l.readPosition >= len(l.input){
		//we set ch to 0 call ASCII 0 is NULL
		l.ch = 0
	}else {
		//if not we set the ch to current position
		l.ch = l.input[l.readPosition]
	}
	//set position to the current position of ch
	l.position = l.readPosition
	//increament the read position by 1
	l.readPosition ++
}

func(l *Lexer)NextToken()token.Token{
	var tok token.Token
	
	l.skipWhiteSpace()
	
	switch l.ch{
		case '=':
			tok = newToken(token.ASSIGN, l.ch)
		case '(':
			tok = newToken(token.LPAREN, l.ch)
		case ')':
			tok = newToken(token.RPAREN, l.ch)
		case '+':
			tok = newToken(token.PLUS, l.ch)
		case '{':
			tok = newToken(token.LBRACE, l.ch)
		case '}':
			tok = newToken(token.RBRACE, l.ch)
		case '-':
			tok = newToken(token.MINUS, l.ch)
		case '!':
			tok = newToken(token.BANG, l.ch)
		case '/':
			tok = newToken(token.SLASH, l.ch)
		case '*':
			tok = newToken(token.ASTERISK, l.ch)
		case '<':
			tok = newToken(token.LT, l.ch)
		case '>':
			tok = newToken(token.GT, l.ch)
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
		case ',':
			tok = newToken(token.COMMA, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
		default:
			if isLetter(l.ch){
				tok.Literal = l.readIndentifier()
				tok.Type = token.LookUpIdent(tok.Literal)
				return tok
			}else if isDigit(l.ch){
				tok.Type = token.INT
				tok.Literal = l.readNumber()
				return tok
			}else {
				tok = newToken(token.ILLEGAL, l.ch)
			}
	}
	
	l.readChar()
	return tok
}

func (l *Lexer)readIndentifier()string{
	position := l.position
	for isLetter(l.ch){
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer)readNumber()string{
	position := l.position
	for isDigit(l.ch){
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer)skipWhiteSpace(){
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

//helpers
func newToken(tokenType token.TokenType, ch byte)token.Token{
	return token.Token{ Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte)bool{
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte)bool{
	return '0' <= ch && ch <= '9'
}

func New(input string)*Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}