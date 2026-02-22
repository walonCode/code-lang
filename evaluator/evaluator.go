package evaluator

import (
	"maps"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/walonCode/code-lang/ast"
	"github.com/walonCode/code-lang/lexer"
	"github.com/walonCode/code-lang/object"
	"github.com/walonCode/code-lang/parser"
)

var moduleCache = map[string]*object.Module{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	//statement
	case *ast.Program:
		return evalProgram(node, *env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.ImportStatement:
		return evalImportStatement(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	//expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.CharLiteral:
		return &object.Char{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node, right)
	case *ast.InfixExpression:
		if isAssignment(node.Operator) {
			return evalAssignment(node, env)
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node, *env)
	case *ast.IfExpression:
		return evalIfExpression(node, *env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: &body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpression(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, node)
	case *ast.ArrayLiteral:
		elements := evalExpression(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index, node)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.MemberExpression:
		obj := Eval(node.Object, env)
		if isError(obj) {
			return obj
		}

		return evalMemberExpression(obj, node)
	}

	return nil
}

func evalImportStatement(node *ast.ImportStatement, env *object.Environment)object.Object{
	modulePath := node.Path
	
	if mod, ok := moduleCache[modulePath];ok {
		return mod
	}
	
	fileName := filepath.Clean(modulePath + ".cl")
	content, err := os.ReadFile(fileName)
	if err != nil {
		return object.NewError(node.Line(), node.Column(), "could not read module %q : %s",modulePath, err)
	}
	
	moduleEnv := object.NewEnclosedEnvironment(env)
	
	l := lexer.New(string(content))
	p := parser.New(l)
	programe := p.ParsePrograme()
	
	Eval(programe, moduleEnv)
	
	moduleobj := &object.Module{ Members: map[string]object.Object{}}
	maps.Copy(moduleobj.Members, moduleEnv.Store)
	env.Set(modulePath, moduleobj)
	
	moduleCache[modulePath] = moduleobj
	
	return moduleobj
}

func isAssignment(op string) bool {
	switch op {
	case "=", "+=", "-=", "*=", "/=", "%=", "**=", "//=":
		return true
	default:
		return false
	}
}

func evalAssignment(node *ast.InfixExpression, env *object.Environment) object.Object {
	val := Eval(node.Right, env)
	if isError(val) {
		return val
	}

	switch left := node.Left.(type) {
	case *ast.Identifier:
		var finalVal object.Object
		if node.Operator == "=" {
			finalVal = val
		} else {
			currentVal := evalIdentifier(left, env)
			if isError(currentVal) {
				return currentVal
			}
			finalVal = evalInfixExpression(node, currentVal, val)
		}

		if isError(finalVal) {
			return finalVal
		}

		_, updated := env.Update(left.Value, finalVal)
		if !updated {
			env.Set(left.Value, finalVal)
		}
		return finalVal

	case *ast.MemberExpression:
		obj := Eval(left.Object, env)
		if isError(obj) {
			return obj
		}

		var finalVal object.Object
		if node.Operator == "=" {
			finalVal = val
		} else {
			currentVal := evalMemberExpression(obj, left)
			if isError(currentVal) {
				return currentVal
			}
			finalVal = evalInfixExpression(node, currentVal, val)
		}

		if isError(finalVal) {
			return finalVal
		}

		return evalAssignMember(obj, left, finalVal)

	case *ast.IndexExpression:
		obj := Eval(left.Left, env)
		if isError(obj) {
			return obj
		}

		idx := Eval(left.Index, env)
		if isError(idx) {
			return idx
		}

		var finalVal object.Object
		if node.Operator == "=" {
			finalVal = val
		} else {
			currentVal := evalIndexExpression(obj, idx, left)
			if isError(currentVal) {
				return currentVal
			}
			finalVal = evalInfixExpression(node, currentVal, val)
		}

		if isError(finalVal) {
			return finalVal
		}

		return evalAssignIndex(obj, idx, finalVal, left)
	}

	return object.NewError(node.Line(), node.Column(), "invalid left-hand side in assignment")
}

func evalAssignIndex(obj, idx, val object.Object, node *ast.IndexExpression) object.Object {
	switch obj := obj.(type) {
	case *object.Array:
		i, ok := idx.(*object.Integer)
		if !ok {
			return object.NewError(node.Line(), node.Column(), "index must be an integer, got %s", idx.Type())
		}
		if i.Value < 0 || i.Value >= int64(len(obj.Elements)) {
			return object.NewError(node.Line(), node.Column(), "index out of range: %d", i.Value)
		}
		obj.Elements[i.Value] = val
		return val
	case *object.Hash:
		hashKey, ok := idx.(object.Hashable)
		if !ok {
			return object.NewError(node.Line(), node.Column(), "unusable as hash key: %s", idx.Type())
		}
		obj.Pairs[hashKey.HashKey()] = object.HashPair{Key: idx, Value: val}
		return val
	default:
		return object.NewError(node.Line(), node.Column(), "index assignment not supported for %s", obj.Type())
	}
}

func evalAssignMember(obj object.Object, node *ast.MemberExpression, val object.Object) object.Object {
	switch obj := obj.(type) {
	case *object.Hash:
		key := &object.String{Value: node.Property.Value}
		obj.Pairs[key.HashKey()] = object.HashPair{
			Key:   key,
			Value: val,
		}

		return val
	case *object.Module:
		obj.Members[node.Property.Value] = val
		return val	
	default:
		return object.NewError(node.Line(), node.Column(), "cannot assign to property %s on %s", node.Property.Value, obj.Type())
	}
}

func evalMemberExpression(obj object.Object, node *ast.MemberExpression) object.Object {
	switch obj := obj.(type) {
	case *object.Hash:
		key := &object.String{Value: node.Property.Value}
		if val, ok := obj.Pairs[object.HashKey(key.HashKey())]; ok {
			return val.Value
		}
		return object.NewError(node.Line(), node.Column(), "property not found: %s", node.Property.Value)
	case *object.Module:
		val, ok := obj.Members[node.Property.Value]
		if !ok {
			return object.NewError(node.Line(), node.Column(), "module has not member %s", node.Property.Value)
		}
		
		return val
	default:
		return object.NewError(node.Line(), node.Column(), "cannot access property %s on %s", node.Property.Value, obj.Type())
	}
}

func evalWhileExpression(node *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object = object.NULL

	for {
		if node.Condition != nil {
			condition := Eval(node.Condition, env)
			if isError(condition) {
				return condition
			}

			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(node.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForExpression(node *ast.ForExpression, env *object.Environment) object.Object {
	forEnv := object.NewEnclosedEnvironment(env)
	if node.Init != nil {
		initRes := Eval(node.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var result object.Object = object.NULL

	for {
		if node.Condition != nil {
			condition := Eval(node.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(node.Body, forEnv)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

		if node.Post != nil {
			postRes := Eval(node.Post, forEnv)
			if isError(postRes) {
				return postRes
			}
		}
	}

	return result
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashkey, ok := key.(object.Hashable)
		if !ok {
			return object.NewError(node.Line(), node.Column(), "unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashkey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index object.Object, node *ast.IndexExpression) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index, node)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return object.NewError(node.Line(), node.Column(), "index operator not supported: %s", left.Type())
	}
}

func evalStringIndexExpression(left, index object.Object) object.Object {
	strObject := left.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value))

	if idx < 0 || idx > max {
		return object.NULL
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func evalHashIndexExpression(hash, index object.Object, node *ast.IndexExpression) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return object.NewError(node.Line(), node.Column(), "unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return object.NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(left, index object.Object) object.Object {
	arrayObj := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return object.NULL
	}

	return arrayObj.Elements[idx]
}

func applyFunction(fn object.Object, args []object.Object, node *ast.CallExpression) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(node, args...)
	default:
		return object.NewError(node.Line(), node.Column(), "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalExpression(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return object.NewError(node.Line(), node.Column(), "identifier not found: %s", node.Value)
}

func evalProgram(program *ast.Program, env object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, &env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatements(block *ast.BlockStatement, env object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, &env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(node *ast.IfExpression, env object.Environment) object.Object {
	condition := Eval(node.Condition, &env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, &env)
	}

	for _, v := range node.IfElse {
		condition := Eval(v.Condition, &env)
		if isError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(v.Consequence, &env)
		}
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, &env)
	}

	return object.NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(node, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatAndIntegerInfixExpression(node, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalFloatAndIntegerInfixExpression(node, left, right)
	case node.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node, left, right)
	case left.Type() == object.CHAR_OBJ && right.Type() == object.CHAR_OBJ:
		return evalCharInfixExpression(node, left, right)
	case isStringOrChar(left) && isStringOrChar(right):
		return evalMixedStringOrCharInfixExpression(node, left, right)
	default:
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalFloatAndIntegerInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	var leftVal float64
	var rightVal float64

	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}

	switch node.Operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f / %f", leftVal, rightVal)
		}
		return &object.Float{Value: leftVal / rightVal}
	case "**":
		return &object.Float{Value: (math.Pow(float64(leftVal), float64(rightVal)))}
	case "//":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f // %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Floor(float64(leftVal) / float64(rightVal)))}
	case "%":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f %% %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Mod(float64(leftVal), float64(rightVal)))}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "/=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f / %f", leftVal, rightVal)
		}
		return &object.Float{Value: leftVal / rightVal}
	case "**=":
		return &object.Float{Value: (math.Pow(float64(leftVal), float64(rightVal)))}
	case "//=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f // %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Floor(float64(leftVal) / float64(rightVal)))}
	case "%=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f %% %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Mod(float64(leftVal), float64(rightVal)))}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalFloatInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch node.Operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f / %f", leftVal, rightVal)
		}
		return &object.Float{Value: leftVal / rightVal}
	case "**":
		return &object.Float{Value: (math.Pow(float64(leftVal), float64(rightVal)))}
	case "//":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f // %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Floor(float64(leftVal) / float64(rightVal)))}
	case "%":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f %% %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Mod(float64(leftVal), float64(rightVal)))}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "/=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f / %f", leftVal, rightVal)
		}
		return &object.Float{Value: leftVal / rightVal}
	case "**=":
		return &object.Float{Value: (math.Pow(float64(leftVal), float64(rightVal)))}
	case "//=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f // %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Floor(float64(leftVal) / float64(rightVal)))}
	case "%=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %f %% %f", leftVal, rightVal)
		}
		return &object.Float{Value: (math.Mod(float64(leftVal), float64(rightVal)))}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s %s %s", strings.ToLower(string(left.Type())), node.Operator, strings.ToLower(string(right.Type())))
	}
}

func evalMixedStringOrCharInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	if node.Operator != "+" {
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Char).Value

	return &object.String{Value: leftVal + string(rightVal)}
}

func evalCharInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Char).Value
	rightVal := right.(*object.Char).Value
	switch node.Operator {
	case "+":
		return &object.String{Value: string(leftVal) + string(rightVal)}
	default:
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalStringInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	if node.Operator != "+" {
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func evalIntegerInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch node.Operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %d / %d", leftVal, rightVal)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "**":
		return &object.Integer{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case "//":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %d // %d", leftVal, rightVal)
		}
		return &object.Integer{Value: int64(math.Floor(float64(leftVal) / float64(rightVal)))}
	case "%":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %d %% %d", leftVal, rightVal)
		}
		return &object.Integer{Value: int64(math.Mod(float64(leftVal), float64(rightVal)))}
	case "+=":
		return &object.Integer{Value: leftVal + rightVal}
	case "-=":
		return &object.Integer{Value: leftVal - rightVal}
	case "*=":
		return &object.Integer{Value: leftVal * rightVal}
	case "/=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %d / %d", leftVal, rightVal)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "**=":
		return &object.Integer{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case "//=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %d // %d", leftVal, rightVal)
		}
		return &object.Integer{Value: int64(math.Floor(float64(leftVal) / float64(rightVal)))}
	case "%=":
		if rightVal == 0 {
			return object.NewError(node.Line(), node.Column(), "division by zero: %d %% %d", leftVal, rightVal)
		}
		return &object.Integer{Value: int64(math.Mod(float64(leftVal), float64(rightVal)))}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalPrefixExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(node, right)
	default:
		return object.NewError(node.Line(), node.Column(), "unknown operator: %s%s", node.Operator, right.Type())
	}
}

func evalMinusOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NewError(node.Line(), node.Column(), "unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isStringOrChar(obj object.Object) bool {
	return obj.Type() == object.STRING_OBJ || obj.Type() == object.CHAR_OBJ
}
