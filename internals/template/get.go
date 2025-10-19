package template

import (
	"fmt"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/path"
)

func templateGet(context context.TemplateContext) any {
	var res string

	templatePath, err := path.GetTemplatePath(context.Path)

	if err != nil {
		fmt.Println("error templating:", err.Error())
		return err.Error()
	}

	templatePath.PathComponent.Value, context = templatePath.PathComponent.ApplyFunc(templatePath.PathComponent.Value, context)
	
	res, context = templatePath.Protocol.ApplyFunc(templatePath.PathComponent.Value, context)

	res, context = templatePath.Operator.ApplyFunc(templatePath.PathComponent.Value, res, context)
	
	res = cleanOutput(res, context)

	return res
}