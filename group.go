package fasthttprouter

import (
	"github.com/valyala/fasthttp"
	"strings"
)

type groupItem struct {
	path    string
	method  string
	handler fasthttp.RequestHandler
}

type Group struct {
	prefix string

	items []groupItem
}

func NewGroup(prefix string) *Group {
	return &Group{
		prefix: prefix,
	}
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (g *Group) GET(path string, handle fasthttp.RequestHandler) {
	g.Handle("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (g *Group) HEAD(path string, handle fasthttp.RequestHandler) {
	g.Handle("HEAD", path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (g *Group) OPTIONS(path string, handle fasthttp.RequestHandler) {
	g.Handle("OPTIONS", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (g *Group) POST(path string, handle fasthttp.RequestHandler) {
	g.Handle("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (g *Group) PUT(path string, handle fasthttp.RequestHandler) {
	g.Handle("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (g *Group) PATCH(path string, handle fasthttp.RequestHandler) {
	g.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (g *Group) DELETE(path string, handle fasthttp.RequestHandler) {
	g.Handle("DELETE", path, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (g *Group) Handle(method, path string, handle fasthttp.RequestHandler) {
	g.items = append(g.items, groupItem{
		method: method,
		path: path,
		handler: handle,
	})
}

func (g *Group) Append(router *Router) {
	for _, v := range g.items {
		var path string

		hasSuf := strings.HasSuffix(g.prefix, "/")
		hasPref := strings.HasPrefix(v.path, "/")

		switch {
		case hasSuf && hasPref:
			path = g.prefix[:len(g.prefix) - 1] + v.path
		case !hasSuf && !hasPref:
			path = g.prefix + "/" + v.path
		default:
			path = g.prefix + v.path
		}

		router.Handle(v.method, path, v.handler)
	}
}
