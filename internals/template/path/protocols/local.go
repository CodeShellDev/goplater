package protocols

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/path/types"
	"github.com/codeshelldev/goplater/utils/fetchutils"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

func init() {
	Register(types.TemplateProtocol{
		Name: "local",
		ApplyFunc: func(pathComponent string, context context.TemplateContext) (string, context.TemplateContext) {
			var res string

			isRel := strings.HasPrefix(pathComponent, "./") || strings.HasPrefix(pathComponent, "../")

			filePathAbs, _ := filepath.Abs(context.Invoker)
			
			relPathComponent := fsutils.Relative(filepath.Dir(filePathAbs), pathComponent)

			context.Invoker = relPathComponent

			if isRel {
				res = fetchutils.Local(pathComponent, filepath.Dir(filePathAbs))
			} else {
				res = fetchutils.Local(pathComponent, context.Options.Source)
			}

			return res, context
		},
		MatchFunc: func(match string) (bool, string) {
			re, err := regexp.Compile(`.*(\S*)(#://)`)

			if err != nil {
				return false, match
			}

			return re.MatchString(match), "#://"
		},
	})
}