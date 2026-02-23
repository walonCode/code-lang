package net

import (
	"fmt"
	"net/http"

	"github.com/walonCode/code-lang/ast"
	"github.com/walonCode/code-lang/object"
)

type ApplyFunctionFunc func(fn object.Object, args []object.Object, node *ast.CallExpression)object.Object

func NetModule(applyFunc ApplyFunctionFunc) *object.Module {
	members := map[string]object.Object{}

	members["server"] = &object.Builtin{
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			return NewServer(applyFunc)
		},
	}

	return &object.Module{
		Members: members,
	}
}

func NewServer(applyFunc ApplyFunctionFunc) object.Object {
	// server := &object.Server{}

	module := &object.Module{
		Members: map[string]object.Object{},
	}

	module.Members["listen"] = &object.Builtin{
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) == 0 || len(args) > 2 {
				return object.NewError(node.Line(), node.Column(), "listen expect 1 arugment")
			}

			port, ok := args[0].(*object.Integer)
			if !ok {
				return object.NewError(node.Line(), node.Column(), "port must be an integer")
			}

			var callback object.Object
			if len(args) == 2 {
				callback = args[1]
			}

			// handler := func(w http.ResponseWriter, r *http.Request) {
			// 	w.Write([]byte("Hello world"))
			// }

			if callback != nil {
				switch cb := callback.(type) {
				case *object.Builtin:
					cb.Fn(nil) // pass nil for node
				case *object.Function:
					applyFunc(cb, []object.Object{}, node)
				default:
					return object.NewError(node.Line(), node.Column(), "callback must be a function")
				}
			}

			err := http.ListenAndServe(fmt.Sprintf(":%d", port.Value), http.HandlerFunc(nil))
			if err != nil {
				fmt.Println("server error", err)
			}

			return nil
		},
	}

	return module
}
