package collections

import (
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/core"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/ext/html"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/ext/markdown"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/ext/yaml"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/container"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/conversion"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/debug"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/encoding"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/http"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/json"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/math"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/regex"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/strings"
)


var Safe = []modules.Module{
	encoding.Module,
	container.Module,
	conversion.Module,
	html.Module,
	json.Module,
	core.Module,
	regex.Module,
	strings.Module,
	yaml.Module,
	math.Module,
	debug.Module,
	markdown.Module,
}

var Unsafe = []modules.Module{
	http.Module,
}