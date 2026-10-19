package funcs

import "github.com/codeshelldev/goplater/pkg/gplt/modules"

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
	walkDirFunc,
	mkdirFunc,
	fsRemoveFunc,

	joinPathFunc,
	basePathFunc,
	pathDirFunc,
	fileExtFunc,

	absPathFunc,
	relPathFunc,
)