package core

import "github.com/codeshelldev/goplater/pkg/templating/modules"

var Module = modules.NewModule(importFunc, callFunc, globalSetFunc, globalGetFunc, getTemplateBodyFunc, getTemplateDataFunc, returnFunc, returnNextFunc, returnAllFunc, returnOutputsFunc, getOutputsFunc)