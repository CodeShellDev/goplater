package collections

import (
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/codeshelldev/goplater/pkg/templating/modules/base64"
	"github.com/codeshelldev/goplater/pkg/templating/modules/container"
	"github.com/codeshelldev/goplater/pkg/templating/modules/conversion"
	"github.com/codeshelldev/goplater/pkg/templating/modules/core"
	"github.com/codeshelldev/goplater/pkg/templating/modules/debug"
	"github.com/codeshelldev/goplater/pkg/templating/modules/html"
	"github.com/codeshelldev/goplater/pkg/templating/modules/http"
	"github.com/codeshelldev/goplater/pkg/templating/modules/json"
	"github.com/codeshelldev/goplater/pkg/templating/modules/math"
	"github.com/codeshelldev/goplater/pkg/templating/modules/regex"
	"github.com/codeshelldev/goplater/pkg/templating/modules/strings"
	"github.com/codeshelldev/goplater/pkg/templating/modules/yaml"
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
}