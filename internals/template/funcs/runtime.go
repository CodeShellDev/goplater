package funcs

type Runtime struct {
    globals map[string]any
    locals  map[string]map[string]any
	funcs   map[string]string
}

var runtime = &Runtime{
	globals: map[string]any{},
	locals: map[string]map[string]any{},
	funcs: map[string]string{},
}