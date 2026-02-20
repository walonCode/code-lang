package ast

import (
	"bytes"
	"strings"

	"github.com/walonCode/code-lang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	Line() int
	Column() int
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
	Name  *Identifier
	Value Expression
}

// method on the let statement
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
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
func (ls *LetStatement) Line() int { return ls.Token.Line }
func (ls *LetStatement) Column() int { return ls.Token.Column }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

// methods on the return statement
func (i *ReturnStatement) statementNode()       {}
func (i *ReturnStatement) TokenLiteral() string { return i.Token.Literal }
func (i *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(i.TokenLiteral() + " ")
	if i.ReturnValue != nil {
		out.WriteString(i.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
func (i *ReturnStatement) Line() int { return i.Token.Line }
func (i *ReturnStatement) Column() int { return i.Token.Column }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// methods on prefix expression
func (i *PrefixExpression) expressionNode()      {}
func (i *PrefixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Operator)
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}
func (i *PrefixExpression) Line() int { return i.Token.Line }
func (i *PrefixExpression) Column() int { return i.Token.Column }

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Right    Expression
	Operator string
}

// method on infix expression
func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}
func (i *InfixExpression) Line() int { return i.Token.Line }
func (i *InfixExpression) Column() int { return i.Token.Column }

type Identifier struct {
	Token token.Token
	Value string
}

// methods on the Identifier
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) Line() int { return i.Token.Line }
func (i *Identifier) Column() int { return i.Token.Column }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// methods on Interger Literal
func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }
func (i *IntegerLiteral) Line() int { return i.Token.Line }
func (i *IntegerLiteral) Column() int { return i.Token.Column }

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// method on the expression statement
func (i *ExpressionStatement) statementNode()       {}
func (i *ExpressionStatement) TokenLiteral() string { return i.Token.Literal }
func (i *ExpressionStatement) String() string {
	if i.Expression != nil {
		return i.Expression.String()
	}
	return ""
}
func (i *ExpressionStatement) Line() int { return i.Token.Line }
func (i *ExpressionStatement) Column() int { return i.Token.Column }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (p *Program) Line() int { 
	if len(p.Statements) > 0 {
		return p.Statements[0].Line()
	}
	return 0
 }
func (p *Program) Column() int { return p.Statements[0].Column() }

// Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
func (b *Boolean) Line() int { return b.Token.Line }
func (b *Boolean) Column() int { return b.Token.Column }

// if expression
type ELSE_IF struct {
	Condition Expression
	Consequence *BlockStatement
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	IfElse []*ELSE_IF
	Alternative *BlockStatement
}

// methods on the if expression
func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
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
func (i *IfExpression) Line() int { return i.Token.Line }
func (i *IfExpression) Column() int { return i.Token.Column }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

// method on the block statement
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (bs *BlockStatement) Line() int { return bs.Token.Line }
func (bs *BlockStatement) Column() int { return bs.Token.Column }

// function
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       BlockStatement
}

// method on the function literal
func (fl *FunctionLiteral) expressionNode()      {}
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
func (fl *FunctionLiteral) Line() int { return fl.Token.Line }
func (fl *FunctionLiteral) Column() int { return fl.Token.Column }

// call expression
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

// methods on the call expression
func (ce *CallExpression) expressionNode()      {}
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
func (ce *CallExpression) Line() int { return ce.Token.Line }
func (ce *CallExpression) Column() int { return ce.Token.Column }

// strings
type StringLiteral struct {
	Token token.Token
	Value string
}

// method on string
func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) Line() int { return sl.Token.Line }
func (sl *StringLiteral) Column() int { return sl.Token.Column }

// char
type CharLiteral struct {
	Token token.Token
	Value rune
}

func (sl *CharLiteral) expressionNode()      {}
func (sl *CharLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *CharLiteral) String() string       { return sl.Token.Literal }
func (sl *CharLiteral) Line() int { return sl.Token.Line }
func (sl *CharLiteral) Column() int { return sl.Token.Column }

// float
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (sl *FloatLiteral) expressionNode()      {}
func (sl *FloatLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *FloatLiteral) String() string       { return sl.Token.Literal }
func (sl *FloatLiteral) Line() int { return sl.Token.Line }
func (sl *FloatLiteral) Column() int { return sl.Token.Column }

// array
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

// method on Array literal
func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (al *ArrayLiteral) Line() int { return al.Token.Line }
func (al *ArrayLiteral) Column() int { return al.Token.Column }

// array index expression
type IndexExpression struct {
	Token token.Token //[
	Left  Expression
	Index Expression
}

// method on array index
func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
func (ie *IndexExpression) Line() int { return ie.Token.Line }
func (ie *IndexExpression) Column() int { return ie.Token.Column }

// hash literal
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
func (hl *HashLiteral) Line() int { return hl.Token.Line }
func (hl *HashLiteral) Column() int { return hl.Token.Column }

