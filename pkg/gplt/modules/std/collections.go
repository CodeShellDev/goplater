package std

import (
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/core"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/container"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/conversion"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/debug"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/encoding"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/hash"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/http"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/json"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/math"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/regex"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/strings"
)

type StandardCollections struct {
	Safe	[]modules.Module
	Unsafe	[]modules.Module
}

var Collections = StandardCollections{
	Safe: []modules.Module{
		encoding.Module,
		hash.Module,
		container.Module,
		conversion.Module,
		json.Module,
		core.Module,
		regex.Module,
		strings.Module,
		math.Module,
		debug.Module,
	},

	Unsafe: []modules.Module{
		http.Module,
	},
}