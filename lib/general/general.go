package general

import (
	"fmt"

	"github.com/walonCode/code-lang/ast"
	"github.com/walonCode/code-lang/object"
)

var GeneralBuiltins = map[string]*object.Builtin{
	"print": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			for _, value := range args {
				fmt.Println(value.Inspect())
			}
			return nil
		},
	},
	"len": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(node.Line(), node.Column(),"wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return object.NewError(node.Line(), node.Column(),"argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
