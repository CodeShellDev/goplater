package templating

type ContextKey string

type Context struct {
	values	map[ContextKey]any
}

func (c *Context) Copy(dst *Context) {
	dst.tryInit()

	DeepCopyContextMap(c.values, dst.values)
}

func (c *Context) tryInit() {
	if c.values == nil {
		c.values = map[ContextKey]any{}
	}
}

func (c *Context) Set(name ContextKey, v any) {
	c.tryInit()

    c.values[name] = v
}

func (c *Context) Get(name ContextKey) any {
	c.tryInit()

    return c.values[name]
}

func (c *Context) Has(name ContextKey) bool {
	c.tryInit()

	_, exists := c.values[name]

	return exists
}

const InputContextKey ContextKey = "templateInput"

type TemplateInputContext struct {
	Data any
	Body string
}

func DeepCopyContextMap(src map[ContextKey]any, dst map[ContextKey]any) map[ContextKey]any {
	if src == nil {
		return nil
	}

	for k, v := range src {
		dst[k] = deepCopyAny(v)
	}

	return dst
}

func deepCopyAny(value any) any {
    switch val := value.(type) {
    case map[string]any:
        copyMap := make(map[string]any, len(val))
		for k, v := range val {
			copyMap[k] = deepCopyAny(v)
		}

        return copyMap
    case []any:
        copySlice := make([]any, len(val))
        for i, s := range val {
            copySlice[i] = deepCopyAny(s)
        }

        return copySlice
    default:
        return val
    }
}
