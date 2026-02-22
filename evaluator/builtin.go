package evaluator

import (
	"github.com/walonCode/code-lang/std/arrays"
	"github.com/walonCode/code-lang/std/general"
)


func init() {
    // preload standard library modules
    moduleCache["arrays"]  = arrays.Module()
    moduleCache["fmt"] = general.Module()
}