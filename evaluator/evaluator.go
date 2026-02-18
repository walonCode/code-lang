package evaluator

import (
	"github.com/walonCode/code-lang/ast"
	"github.com/walonCode/code-lang/object"
)

func Eval(node ast.Node)object.Object{
	switch node := node.(type){
		//statement
		case *ast.Program:
			return evalStatements(node.Statements)
		case *ast.ExpressionStatement:
			return Eval(node.Expression)
		
		//expression
		case *ast.IntegerLiteral:
			return &object.Integer{ Value: node.Value}
	}
	
	return nil
}

func evalStatements(stmt []ast.Statement)object.Object{
	var result object.Object
	
	for _, statement := range stmt {
		result = Eval(statement)
	}
	
	return result
}