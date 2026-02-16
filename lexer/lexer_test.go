package lexer

import (
	"testing"

	"github.com/walonCode/code-lang/token"
)

func TestNextToken(t *testing.T){
	input := `=+(){},;`
	
	test := []struct{
		expectedType token.TokenType
		expectedLiteral string
	}{
		
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	
	l := New(input)
	
	for i, tt := range test {
		tok := l.NextToken()
		
		if tok.Type != tt.expectedType{
			t.Fatalf("tests[%d] - tokentype wrong. expected=%v, got=%v",
				i, tt.expectedType,tok.Type)
		}
		
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
			i, tt.expectedLiteral, tok.Literal)
		}
	}
}