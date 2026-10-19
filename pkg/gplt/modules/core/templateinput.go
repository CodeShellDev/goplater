package core

import (
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var getTemplateBodyFunc = modules.NewFunc("getTemplateBody", getTemplateBody)

func getTemplateBody(_ *gplt.Runtime, context *gplt.Context) string {
	ctx := context.Get(gplt.InputContextKey).(gplt.TemplateInputContext)

	return ctx.Body
}

var getTemplateDataFunc = modules.NewFunc("getTemplateData", getTemplateData)

func getTemplateData(_ *gplt.Runtime, context *gplt.Context) any {
	ctx := context.Get(gplt.InputContextKey).(gplt.TemplateInputContext)

	return ctx.Data
}