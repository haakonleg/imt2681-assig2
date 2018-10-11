package router

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(*Request)

// The routes are stored in a trie structure
type routeNode struct {
	children map[string]*routeNode
	handlers map[string]HandlerFunc
}

func (rt *routeNode) Print() {
	i := 0
	tree := rt.children
	for len(tree) != 0 {
		i++
		for k, v := range tree {
			fmt.Printf("Depth%d %s \n", i, k)
			tree = v.children
		}
	}
	fmt.Printf("\n\n\n")
}

func (rt *routeNode) addRoute(method string, path string, handler HandlerFunc) {
	subpaths := strings.Split(path, "/")

	currTree := rt
	for _, p := range subpaths {
		if len(p) == 0 {
			continue
		}

		// This route contains a variable
		if p[0] == '{' && p[len(p)-1] == '}' {
			p = "{var}"
		}

		tree, ok := currTree.children[p]
		if !ok {
			tree = &routeNode{
				children: make(map[string]*routeNode, 0),
				handlers: make(map[string]HandlerFunc, 0)}
			currTree.children[p] = tree
		}
		currTree = tree
	}

	currTree.handlers[method] = handler
}

func (rt *routeNode) resolveRoute(method string, path string) (*routeNode, []string) {
	subpaths := strings.Split(path, "/")

	vars := make([]string, 0)
	currNode := rt
	for _, p := range subpaths {
		if len(p) == 0 {
			continue
		}

		// First try to match it against a constant route
		tree, ok := currNode.children[p]
		if !ok {
			// Find the variable route
			tree, ok = currNode.children["{var}"]
			if ok {
				vars = append(vars, p)
			}
		}
		currNode = tree
	}

	return currNode, vars
}

type Router struct {
	routes routeNode
}

func NewRouter() Router {
	return Router{
		routes: routeNode{
			children: make(map[string]*routeNode, 0),
			handlers: make(map[string]HandlerFunc, 0)}}
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Find the handler for this path
	route, vars := ro.routes.resolveRoute(r.Method, r.URL.Path)
	if route != nil {
		// Check if there is a handler for the HTTP method
		handler, ok := route.handlers[r.Method]
		if ok {
			// Call the handler
			handler(&Request{w, r, vars})
			return
		}
	}
	http.NotFound(w, r)
}

func (ro *Router) Handle(method string, path string, handler HandlerFunc) {
	ro.routes.addRoute(method, path, handler)
}
