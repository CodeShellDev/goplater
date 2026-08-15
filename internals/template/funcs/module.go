package funcs

import "github.com/codeshelldev/goplater/pkg/templating/modules"

var Module = modules.NewModule(
	readFunc, 
	readRawFunc, 
	readArgsFunc, 

	outputToFunc,

	writeFunc,
	appendFileFunc,

	fsExistsFunc,
	isFileFunc,
	isDirFunc,

	listDirFunc,
	mkdirFunc,
	fsRemoveFunc,

	joinPathFunc,
	basePathFunc,
	pathDirFunc,
	fileExtFunc,

	absPathFunc,
	relPathFunc,
)