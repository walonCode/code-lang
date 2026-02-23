package evaluator

import (
	"github.com/walonCode/code-lang/std/arrays"
	"github.com/walonCode/code-lang/std/general"
	"github.com/walonCode/code-lang/std/json"
	"github.com/walonCode/code-lang/std/net"
)


func init() {
    // preload standard library modules
    moduleCache["arrays"]  = arrays.Module()
    moduleCache["fmt"] = general.Module()
    moduleCache["http"] = net.HttpModule()
    moduleCache["json"] = json.JsonModule()
    moduleCache["net"] = net.NetModule()
}