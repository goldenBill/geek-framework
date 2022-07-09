package geekweb

import (
	"fmt"
	"net/http"
)

// HanlderFunc defines the request handler used by geekweb
type HanlderFunc func(w http.ResponseWriter, r *http.Request)

// Engine implements the interface of ServeHTTP
type Engine struct {
	router map[string]HanlderFunc
}

// New is the constructor of geekweb.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HanlderFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HanlderFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HanlderFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HanlderFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
