package core

import "github.com/codeshelldev/goplater/pkg/gplt/modules"

var Module = modules.NewModule(
	importFunc, 
	importExactFunc, 

	callFunc, 

	globalSetFunc, 
	globalGetFunc, 
	globalHasFunc, 

	getTemplateBodyFunc, 
	getTemplateDataFunc, 

	returnFunc, 
	returnNextFunc, 
	returnAllFunc, 
	getOutputsFunc,
)