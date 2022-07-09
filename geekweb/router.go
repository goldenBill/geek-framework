package geekweb

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern only allows one *
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	parts := parsePattern(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(parts, 0)
	if n == nil {
		return nil, nil
	}

	params := make(map[string]string)
	matchedParts := parsePattern(n.pattern)
	for index, part := range matchedParts {
		if part[0] == ':' {
			params[part[1:]] = parts[index]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(parts[index:], "/")
			break
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s for URL %q\n", c.Method, c.Path)
		return
	}
	c.Params = params
	key := c.Method + "-" + n.pattern
	r.handlers[key](c)
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}
