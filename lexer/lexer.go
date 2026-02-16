package lexer

import "github.com/walonCode/code-lang/token"

type Lexer struct {
	input string
	positon int
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
	l.positon = l.readPosition
	//increament the read position by 1
	l.readPosition ++
}

func(l *Lexer)NextToken()token.Token{
	var tok token.Token
	
	switch l.ch{
		case '=':
			tok = newToken(token.ASSIGN, l.ch)
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
		case '(':
			tok = newToken(token.LPAREN, l.ch)
		case ')':
			tok = newToken(token.RPAREN, l.ch)
		case ',':
			tok = newToken(token.COMMA, l.ch)
		case '+':
			tok = newToken(token.PLUS, l.ch)
		case '{':
			tok = newToken(token.LBRACE, l.ch)
		case '}':
			tok = newToken(token.RBRACE, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
	}
	l.readChar()
	return tok
}

//helpers
func newToken(tokenType token.TokenType, ch byte)token.Token{
	return token.Token{ Type: tokenType, Literal: string(ch)}
}

func New(input string)*Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}