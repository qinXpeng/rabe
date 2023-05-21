package rabes

type router struct {
	handlers map[string][]HandlerCtl
}

func FormatRouterKey(method string, path string) string {
	return method + "|" + path
}

func newRouter() *router {
	r := &router{
		handlers: make(map[string][]HandlerCtl),
	}
	return r
}

func (r *router) addRoute(method string, path string, handler ...HandlerCtl) {
	key := FormatRouterKey(method, path)
	r.handlers[key] = handler
}

func (r *router) handle(ctx *Context) {
	key := FormatRouterKey(ctx.Method, ctx.Path)
	if handler, ok := r.handlers[key]; ok {
		ctx.middlewares = append(ctx.middlewares, handler...)
	} else {
		ctx.String(STATUS_NOT_FOUND.Code(), "%s: %s\n", STATUS_NOT_FOUND.Error(), ctx.Path)
	}
	ctx.Next()
}
