package evaluator

import (
	"fmt"

	"github.com/procode2/lena/ln_parser/ast"
	"github.com/procode2/lena/ln_parser/object"
)

type JavaScriptCode struct {
	Code string
}

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		// return NULL
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -(value)}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		// return NULL
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	default:
		// return NULL
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}
func Eval(node ast.Node) JavaScriptCode {
	//fmt.Println("node")
	//fmt.Println(node)
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		// Check for "puts" and replace with "console.log"
		if callExp, ok := node.Expression.(*ast.CallExpression); ok {
			if ident, ok := callExp.Function.(*ast.Identifier); ok {
				if ident.Value == "puts" {
					// Replace "puts" with "console.log"
					callExp.Function = &ast.Identifier{Value: "console.log"}
				}
			}
		}
		return Eval(node.Expression)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		operator := node.Operator
		return JavaScriptCode{Code: fmt.Sprintf("(%s %s %s)", left.Code, operator, right.Code)}
	case *ast.IntegerLiteral:
		return JavaScriptCode{Code: fmt.Sprintf("%d", node.Value)}
	case *ast.Boolean:
		if node.Value {
			return JavaScriptCode{Code: "true"}
		}
		return JavaScriptCode{Code: "false"}
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		operator := node.Operator
		return JavaScriptCode{Code: fmt.Sprintf("(%s%s)", operator, right.Code)}
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		return JavaScriptCode{Code: fmt.Sprintf("return %s;", Eval(node.ReturnValue).Code)}
	case *ast.LetStatement:
		return JavaScriptCode{Code: fmt.Sprintf("let %s = %s;", node.Name.Value, Eval(node.Value).Code)}
	case *ast.Identifier:
		return JavaScriptCode{Code: node.Value}
	case *ast.FunctionLiteral:
		parameters := parametersToString(node.Parameters)
		body := Eval(node.Body)
		return JavaScriptCode{Code: fmt.Sprintf("function(%s) {\n%s\n}", parameters, body.Code)}
	case *ast.CallExpression:
		function := Eval(node.Function)
		arguments := argumentsToString(node.Arguments)
		return JavaScriptCode{Code: fmt.Sprintf("%s(%s)", function.Code, arguments)}
	case *ast.ArrayLiteral:
		elements := elementsToString(node.Elements)
		return JavaScriptCode{Code: fmt.Sprintf("[%s]", elements)}
	case *ast.IndexExpression:
		left := Eval(node.Left)
		index := Eval(node.Index)
		return JavaScriptCode{Code: fmt.Sprintf("%s[%s]", left.Code, index.Code)}
	case *ast.HashLiteral:
		pairs := hashPairsToString(node.Pairs)
		return JavaScriptCode{Code: fmt.Sprintf("{%s}", pairs)}
	case *ast.StringLiteral:
		return JavaScriptCode{Code: fmt.Sprintf(`"%s"`, node.Value)}
	}

	return JavaScriptCode{}
}

func evalProgram(program *ast.Program) JavaScriptCode {
	var accumulatedCode string

	for _, stmt := range program.Statements {
		result := Eval(stmt)
		accumulatedCode += result.Code + "\n"
	}

	return JavaScriptCode{Code: accumulatedCode}
}

func evalBlockStatement(block *ast.BlockStatement) JavaScriptCode {
	var result JavaScriptCode

	for _, stmt := range block.Statements {
		result = Eval(stmt)
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) JavaScriptCode {
	condition := Eval(ie.Condition)
	consequence := Eval(ie.Consequence)
	if ie.Alternative != nil {
		alternative := Eval(ie.Alternative)
		return JavaScriptCode{Code: fmt.Sprintf("if (%s) {\n%s\n} else {\n%s\n}", condition.Code, consequence.Code, alternative.Code)}
	}

	return JavaScriptCode{Code: fmt.Sprintf("if (%s) {\n%s\n}", condition.Code, consequence.Code)}
}

func parametersToString(parameters []*ast.Identifier) string {
	var result string
	for i, param := range parameters {
		result += param.Value
		if i < len(parameters)-1 {
			result += ", "
		}
	}
	return result
}

func argumentsToString(arguments []ast.Expression) string {
	var result string
	for i, arg := range arguments {
		result += Eval(arg).Code
		if i < len(arguments)-1 {
			result += ", "
		}
	}
	return result
}

func elementsToString(elements []ast.Expression) string {
	var result string
	for i, element := range elements {
		result += Eval(element).Code
		if i < len(elements)-1 {
			result += ", "
		}
	}
	return result
}

func hashPairsToString(pairs map[ast.Expression]ast.Expression) string {
	var result string
	for key, value := range pairs {
		result += fmt.Sprintf("%s: %s", Eval(key).Code, Eval(value).Code)
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
