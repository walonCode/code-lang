package ast

import (
	"bytes"
	"strings"

	"github.com/walonCode/code-lang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node // TokenLiteral and String()
	expressionNode()
}

type Program struct {
	Statements []Statement 
}

type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}
//method on the let statement 
func(ls *LetStatement)statementNode(){}
func(ls *LetStatement)TokenLiteral()string{ return ls.Token.Literal}
func(ls *LetStatement)String()string{
	var out bytes.Buffer
	
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	
	out.WriteString(";")
	
	return out.String()
}

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}
//methods on the return statement
func(i *ReturnStatement)statementNode(){}
func (i *ReturnStatement)TokenLiteral()string { return i.Token.Literal}
func(i *ReturnStatement)String()string{
	var out bytes.Buffer
	out.WriteString(i.TokenLiteral() + " ")
	if i.ReturnValue != nil {
		out.WriteString(i.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}
//methods on prefix expression
func(i *PrefixExpression)expressionNode(){}
func (i *PrefixExpression)TokenLiteral()string { return i.Token.Literal}
func (i *PrefixExpression)String()string {
	var out bytes.Buffer
	
	out.WriteString("(")
	out.WriteString(i.Operator)
	out.WriteString(i.Right.String())
	out.WriteString(")")
	
	return out.String()
}

type InfixExpression struct {
	Token token.Token
	Left Expression
	Right Expression
	Operator string
}
//method on infix expression
func(i *InfixExpression)expressionNode(){}
func (i *InfixExpression)TokenLiteral()string { return i.Token.Literal}
func (i *InfixExpression)String()string {
	var out bytes.Buffer
	
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	
	return out.String()
}


type Identifier struct {
	Token token.Token
	Value string
}
//methods on the Identifier
func(i *Identifier)expressionNode(){}
func (i *Identifier)TokenLiteral()string { return i.Token.Literal}
func (i *Identifier)String()string { return i.Value}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}
//methods on Interger Literal
func(i *IntegerLiteral)expressionNode(){}
func (i *IntegerLiteral)TokenLiteral()string { return i.Token.Literal}
func (i *IntegerLiteral)String()string { return i.Token.Literal}

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}
//method on the expression statement
func(i *ExpressionStatement)statementNode(){}
func(i *ExpressionStatement)TokenLiteral()string { return i.Token.Literal }
func(i *ExpressionStatement)String()string{
	if i.Expression != nil {
		return i.Expression.String()
	}
	return ""
}

func(p *Program)TokenLiteral()string{
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}else{
		return ""
	}
}

func (p *Program)String()string{
	var out bytes.Buffer
	
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	
	return out.String()
}

//Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func(b *Boolean)expressionNode(){}
func(b *Boolean)TokenLiteral()string { return b.Token.Literal}
func(b *Boolean)String()string { return b.Token.Literal }


//if expression
type IfExpression struct {
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}
//methods on the if expression
func(i *IfExpression)expressionNode(){}
func(i *IfExpression)TokenLiteral()string { return i.Token.Literal}
func(i *IfExpression)String()string { 
	var out bytes.Buffer
	
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	
	if i.Alternative != nil {
		out.WriteString("else")
		out.WriteString(i.Alternative.String())
	}
	
	return out.String()
}

type BlockStatement struct {
	Token token.Token
	Statements []Statement
}
//method on the block statement
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	
	return out.String()
}

//function
type FunctionLiteral struct {
	Token token.Token
	Parameters []*Identifier
	Body BlockStatement
}
//method on the function literal
func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

//call expression
type CallExpression struct {
	Token token.Token
	Function Expression
	Arguments []Expression
}
//methods on the call expression
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}