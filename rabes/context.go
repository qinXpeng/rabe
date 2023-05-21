package rabes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type H map[string]any

type Context struct {
	Writer http.ResponseWriter

	Req http.Request

	Path   string
	Method string

	Params      map[string]string
	StatusCode  int
	middlewares []HandlerCtl
	index       int
	keys        map[string]any
	mu          sync.RWMutex
}

//
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    *r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
		mu:     sync.RWMutex{},
		keys:   make(map[string]any),
	}

}

func (ctx *Context) Next() {
	ctx.index++
	s := len(ctx.middlewares)
	for ctx.index < s {
		ctx.middlewares[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *Context) Set(key string, value any) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.keys[key] = value
}

func (ctx *Context) Get(key string) (value any, exists bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	value, exists = ctx.keys[key]
	return
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.Req.Header.Get(key)
}

func (ctx *Context) String(code int, format string, values ...any) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	if len(values) == 0 {
		ctx.Writer.Write([]byte(format))
		return
	}
	res := fmt.Sprintf(format, values...)
	ctx.Writer.Write([]byte(res))
}

func (ctx *Context) JSON(code int, obj any) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}
