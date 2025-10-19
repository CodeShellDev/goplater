package path

import (
	"cmp"
	"errors"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/codeshelldev/goplater/internals/template/path/operators"
	"github.com/codeshelldev/goplater/internals/template/path/pathcomponents"
	"github.com/codeshelldev/goplater/internals/template/path/protocols"
	"github.com/codeshelldev/goplater/internals/template/path/types"
)

func GetRegex() string {
	return `((\S)?(\S+/)([^\t\n\r\f\v{}]+[^\s{}]))`
}

func GetTemplatePath(path string) (types.TemplatePath, error) {
	re, err := regexp.Compile(GetRegex())

	if err != nil {
		return types.TemplatePath{}, err
	}

	matches := re.FindStringSubmatch(path)

	if len(matches) < 1 {
		return types.TemplatePath{}, errors.New("invalid syntax")
	}
	
	protocol := getProtocolByRegex(matches[1])

	if protocol.ApplyFunc == nil {
		return types.TemplatePath{}, errors.New("invalid protocol")
	}

	parts := strings.Split(path, protocol.Value)

	operator := getOperatorByRegex(parts[0], protocol)

	if operator.ApplyFunc == nil {
		return types.TemplatePath{}, errors.New("operator does not apply")
	}

	pathComponent := getPathComponentByRegex(parts[1], protocol)

	if pathComponent.ApplyFunc == nil {
		return types.TemplatePath{}, errors.New("unsupported path component")
	}

	return types.TemplatePath{
		Operator: operator,
		Protocol: protocol,
		PathComponent: pathComponent,
	}, nil
}

func getPathComponentByRegex(str string, prot types.TemplateProtocol) types.TemplatePathComponent {
	ordered := slices.Collect(maps.Values(pathcomponents.All()))

	slices.SortFunc(ordered, func(a, b types.TemplatePathComponent) int {
		return cmp.Compare(b.Priority, a.Priority)
	})

	for _, pathComp := range ordered {
		ok, new := pathComp.MatchFunc(str)

		if ok && slices.ContainsFunc(pathComp.AllowedProtocols, func(try types.TemplateProtocol) bool {
			return try.Name == prot.Name
		}) {
			pathComp.Value = new
			return pathComp
		}
	}

	return types.TemplatePathComponent{}
}

func getProtocolByRegex(str string) types.TemplateProtocol {
	ordered := slices.Collect(maps.Values(protocols.All()))

	slices.SortFunc(ordered, func(a, b types.TemplateProtocol) int {
		return cmp.Compare(b.Priority, a.Priority)
	})

	for _, protocol := range ordered {
		ok, new := protocol.MatchFunc(str)

		if ok {
			protocol.Value = new
			return protocol
		}
	}

	return types.TemplateProtocol{}
}

func getOperatorByRegex(str string, prot types.TemplateProtocol) types.TemplateOperator {
	ordered := slices.Collect(maps.Values(operators.All()))

	slices.SortFunc(ordered, func(a, b types.TemplateOperator) int {
		return cmp.Compare(b.Priority, a.Priority)
	})

	for _, operator := range ordered {
		ok, new := operator.MatchFunc(str)

		if ok && slices.ContainsFunc(operator.AllowedProtocols, func(try types.TemplateProtocol) bool {
			return try.Name == prot.Name
		}) {
			operator.Value = new
			return operator
		}
	}

	return types.TemplateOperator{}
}