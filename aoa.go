package aoa

import (
	"net/http"
	"regexp"
)

/**
	比较简单的go http 中间件实现
**/

//Next ...
type Next func()

//HTTPContext ...
type HTTPContext struct {
	req *http.Request
	res http.ResponseWriter
}

type route struct {
	pattern *regexp.Regexp
	handler Handler
}

//Handler ...
type Handler interface {
	ServeHTTP(HTTPContext)
}

//A Middleware type
type Middleware func(HTTPContext, Next)

// App applican
type App struct {
	middlewares []Middleware
	index       int
	ctx         HTTPContext
	mlen        int
	routes      []*route
}

//HandlerFunc ...
type HandlerFunc func(HTTPContext)

func (f HandlerFunc) ServeHTTP(ctx HTTPContext) {
	f(ctx)
}

//HandleFunc ...
func (a *App) HandleFunc(pattern *regexp.Regexp, handler func(HTTPContext)) {
	a.Handler(pattern, HandlerFunc(handler))
}

//Handler ...
func (a *App) Handler(pattern *regexp.Regexp, handler HandlerFunc) {
	a.routes = append(a.routes, &route{pattern, handler})
}

// MiddlewareFunc is middleware use func
func (a *App) MiddlewareFunc(m Middleware) {
	a.middlewares = append(a.middlewares, m)
}

func (a *App) compose(ctx HTTPContext) {
	a.index = 0
	a.mlen = len(a.middlewares)
	if a.mlen <= a.index {
		a.router()
	} else {
		a.middlewares[a.index](a.ctx, a.next)
	}
}

func (a *App) next() {
	a.index++
	if a.index < a.mlen {
		a.middlewares[a.index](a.ctx, a.next)
	} else {
		a.router()
	}
}

func (a *App) router() {
	for _, route := range a.routes {
		if route.pattern.MatchString(a.ctx.req.URL.Path) {
			route.handler.ServeHTTP(a.ctx)
			return
		}
	}
	http.NotFound(a.ctx.res, a.ctx.req)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.ctx = HTTPContext{res: w, req: r}
	a.compose(a.ctx)
}

//NewAppServe is return app examples
func NewAppServe() *App {
	return &App{}
}

//ListenServe is http server listen address and run server
func (a *App) ListenServe(address string) {
	http.ListenAndServe(address, a)
}
