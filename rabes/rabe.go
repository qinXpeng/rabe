package rabes

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerCtl func(c *Context)

type Engine struct {
	*RouteGroup
	router *router
	groups []*RouteGroup
}

type RouteGroup struct {
	prefix      string
	engine      *Engine
	middlewares []HandlerCtl
	preGroup    *RouteGroup
}

func New() *Engine {
	e := &Engine{router: newRouter()}
	e.RouteGroup = &RouteGroup{engine: e}
	e.groups = []*RouteGroup{e.RouteGroup}
	return e
}

func (e *Engine) addRoute(method string, path string, handler ...HandlerCtl) {
	e.router.addRoute(method, path, handler...)
}

func (e *Engine) GET(path string, ctl ...HandlerCtl) {
	e.addRoute("GET", path, ctl...)
}

func (e *Engine) POST(path string, ctl ...HandlerCtl) {
	e.addRoute("POST", path, ctl...)
}

func (e *Engine) Run(addr string) (err error) {

	fmt.Println("Server is running at " + addr + " ...")

	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerCtl
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, r)
	c.middlewares = middlewares
	e.router.handle(c)
}

func (g *RouteGroup) Group(prefix string) *RouteGroup {
	e := g.engine
	newGroup := &RouteGroup{
		prefix:   g.prefix + prefix,
		engine:   e,
		preGroup: g,
	}
	e.groups = append(e.groups, newGroup)
	return newGroup
}

func (g *RouteGroup) GET(path string, ctl ...HandlerCtl) {
	g.engine.GET(g.prefix+path, ctl...)
}
func (g *RouteGroup) POST(path string, ctl ...HandlerCtl) {
	g.engine.POST(g.prefix+path, ctl...)
}

func (g *RouteGroup) Use(middlewares ...HandlerCtl) {
	g.middlewares = append(g.middlewares, middlewares...)
}
