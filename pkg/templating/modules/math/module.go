package math

import "github.com/codeshelldev/goplater/pkg/templating/modules"

var Module = modules.NewModule(
	addFunc, 
	subFunc, 
	multFunc, 
	divdFunc, 
	modFunc,

	absFunc,
	minFunc,
	maxFunc,
	powFunc,
	roundFunc,
	floorFunc,
	ceilFunc,
	sqrtFunc,
	clampFunc,
)
