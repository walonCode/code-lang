package parser

import (
	"testing"
	"fmt"

	"github.com/walonCode/code-lang/ast"
	"github.com/walonCode/code-lang/lexer"
)

func TestLetStatments(t *testing.T){
	input := `
	let x = y;
	let y = 10;
	let foobar = 838383;
	`
	
	l := lexer.New(input)
	p := New(l)
	
	program := p.ParsePrograme()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}
	
	tests := []struct{
		expectedIdentifier string 
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	
	for i,tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier){
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser){
	errors := p.Errors()
	
	if len(errors) == 0 {
		return 
	}
	
	t.Errorf("parser has %d errors,\n", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q\n",msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string)bool{
	if s.TokenLiteral() != "let"{
		t.Errorf("s.TokenLiteral() not 'let'. got=%v\n",s.TokenLiteral())
		return false
	}
	
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatment, got=%T\n",s)
		return false
	}
	
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	
	return true
}

func TestReturnStatement(t *testing.T){
	input := `
	return 5;
	return 10;
	return rrrrrr;
	`
	
	l := lexer.New(input)
	p := New(l)
	
	program := p.ParsePrograme()
	checkParserErrors(t,p)
	
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}
	
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}


func TestIdentifierExpression(t *testing.T){
	input := "foobar;"
	
	l := lexer.New(input)
	p := New(l)
	
	program := p.ParsePrograme()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements, got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement, got=%T",program.Statements[0])
	}
	
	ident,ok := stmt.Expression.(*ast.Indentifier)
	if !ok{
		t.Fatalf("exp not *ast.Indentifier. got=%T", stmt.Expression)
	}
	
	if ident.Value != "foobar" {
		t.Errorf("ident value not %s. got=%s", "foobar", ident.Value)
	}
	
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}


func TestIntergerLiteralExpression(t *testing.T){
	input := "5;"
	
	l := lexer.New(input)
	p := New(l)
	
	program := p.ParsePrograme()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements, got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement, got=%T",program.Statements[0])
	}
	
	ident,ok := stmt.Expression.(*ast.IntergerLiteral)
	if !ok{
		t.Fatalf("exp not *ast.Indentifier. got=%T", stmt.Expression)
	}
	
	if ident.Value != 5 {
		t.Errorf("ident value not %d. got=%d", 5, ident.Value)
	}
	
	if ident.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", "5", ident.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T){
	prefixTest := []struct{
		input string
		operator string
		intergerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-",15},
		{"!4", "!", 4},
	}
	
	for _,tt := range prefixTest {
		l := lexer.New(tt.input)
		p := New(l)
		
		program := p.ParsePrograme()
		checkParserErrors(t,p)
		
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}
		
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		
		if !testIntegerLiteral(t, exp.Right, tt.intergerValue) {
			return
		}
	}
}


func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntergerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
	return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
	return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}


func TestParsingInfixExpression(t *testing.T){
	infixTest := []struct {
		input string
		leftValue int64
		operator string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	
	for _, tt := range infixTest {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParsePrograme()
		checkParserErrors(t,p)
		
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}
		
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		expected string
	}{
	{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParsePrograme()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}