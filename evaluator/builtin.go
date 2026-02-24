package evaluator

import (
	"github.com/walonCode/code-lang/std/arrays"
	"github.com/walonCode/code-lang/std/fs"
	"github.com/walonCode/code-lang/std/general"
	"github.com/walonCode/code-lang/std/json"
	"github.com/walonCode/code-lang/std/math"
	"github.com/walonCode/code-lang/std/net"
	"github.com/walonCode/code-lang/std/strings"
	"github.com/walonCode/code-lang/std/time"
)

func init() {
	// preload standard library modules
	moduleCache["arrays"] = arrays.Module()
	moduleCache["fmt"] = general.Module()
	moduleCache["http"] = net.HttpModule()
	moduleCache["json"] = json.JsonModule()
	moduleCache["net"] = net.NetModule(applyFunction)
	moduleCache["fs"] = fs.Module()
	moduleCache["math"] = math.Module()
	moduleCache["strings"] = strings.Module()
    moduleCache["time"] = time.Module()
}
