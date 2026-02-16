package ast

import "github.com/walonCode/code-lang/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement 
}

type LetStatement struct {
	Token token.Token
	Name *Indentifier
	Value Expression
}

func(ls *LetStatement)statementNode(){}
func(ls *LetStatement)TokenLiteral()string{ return ls.Token.Literal}

type Indentifier struct {
	Token token.Token
	Value string
}

func(p *Program)TokenLiteral()string{
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}else{
		return ""
	}
}