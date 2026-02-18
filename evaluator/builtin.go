package evaluator

import (
	"maps"

	"github.com/walonCode/code-lang/lib/arrays"
	"github.com/walonCode/code-lang/lib/general"
	"github.com/walonCode/code-lang/object"
)

var builtins = map[string]*object.Builtin{}

func init() {
	maps.Copy(builtins, general.GeneralBuiltins)
	maps.Copy(builtins, arrays.ArrayBuiltins)
}
