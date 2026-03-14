package core

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var getTemplateBodyFunc = modules.NewFunc("getTemplateBody", getTemplateBody)

func getTemplateBody(_ *templating.Runtime, context templating.Context) string {
	ctx := context.Get(templating.InputContextKey).(templating.TemplateInputContext)

	return ctx.Body
}

var getTemplateDataFunc = modules.NewFunc("getTemplateData", getTemplateData)

func getTemplateData(_ *templating.Runtime, context templating.Context) any {
	ctx := context.Get(templating.InputContextKey).(templating.TemplateInputContext)

	return ctx.Data
}