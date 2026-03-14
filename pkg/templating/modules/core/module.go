package core

import "github.com/codeshelldev/goplater/pkg/templating/modules"

var Module = modules.NewModule(funcCallFunc, funcCallWithArgsFunc, funcDefineFunc, globalSetFunc, globalGetFunc, getTemplateBodyFunc, getTemplateDataFunc)