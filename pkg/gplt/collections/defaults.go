package collections

import (
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/base64"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/container"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/conversion"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/core"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/debug"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/html"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/http"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/json"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/markdown"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/math"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/regex"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/strings"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/yaml"
)


var All = []modules.Module{
	base64.Module,
	container.Module,
	conversion.Module,
	html.Module,
	http.Module,
	json.Module,
	core.Module,
	regex.Module,
	strings.Module,
	yaml.Module,
	math.Module,
	debug.Module,
	markdown.Module,
}