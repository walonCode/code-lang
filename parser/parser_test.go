package parser

import (
	"testing"

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