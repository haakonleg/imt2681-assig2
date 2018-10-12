/*
	Package router implements a simple router that can route requests based on the URL paths/subpaths and
	HTTP verbs. Dynamic "variables" are supported, by enclosing a path name in curly brackets. Validator
	functions for these can also be registered, where if false is returned, a 404 status code will be sent.
*/

package router

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc is the function template for handlers
type HandlerFunc func(*Request)

// ValidatorFunc is the function template for validators
// It recieves the variable as a string, and returns a bool indicating if the validation was successful
// and a generic interface{} where variables can be decoded to other types
type ValidatorFunc func(string) (bool, interface{})

// The routes are stored in a trie structure
type routeNode struct {
	children map[string]*routeNode
	varNames map[string]string
	handlers map[string]HandlerFunc
}

// Print prints the routeNode tree
func (rt *routeNode) Print(depth int) {
	tree := rt.children
	for k, v := range tree {
		fmt.Printf("Depth %d %s \n", depth, k)
		fmt.Printf("Vars: %v, Handlers: %v\n", v.varNames, v.handlers)
		if v.children != nil {
			v.Print(depth + 1)
		}
	}
}

func (rt *routeNode) addRoute(method string, path string, handler HandlerFunc) {
	subpaths := strings.Split(path, "/")

	varName := ""
	currNode := rt
	for _, p := range subpaths {
		if len(p) == 0 {
			continue
		}

		// This route contains a variable
		if p[0] == '{' && p[len(p)-1] == '}' {
			varName = p[1 : len(p)-1]
			p = "{var}"
		}

		node, ok := currNode.children[p]
		if !ok {
			node = &routeNode{
				children: make(map[string]*routeNode, 0),
				varNames: make(map[string]string, 0),
				handlers: make(map[string]HandlerFunc, 0)}
			currNode.children[p] = node
		}
		currNode = node
	}

	if varName != "" {
		currNode.varNames[method] = varName
	}
	currNode.handlers[method] = handler
}

// Walks the route tree and finds a routeNode corresponding to the specified method and URL path
// on the way it also captures all variables into a slice
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
			// Try to find a variable route
			tree, ok = currNode.children["{var}"]
			if ok {
				vars = append(vars, tree.varNames[method])
				vars = append(vars, p)
			}
		}
		currNode = tree
	}

	return currNode, vars
}

// Router is the context for a router object
type Router struct {
	routes     routeNode
	validators map[string]ValidatorFunc
}

// NewRouter creates a new Router object
func NewRouter() *Router {
	return &Router{
		routes: routeNode{
			children: make(map[string]*routeNode, 0),
			varNames: make(map[string]string, 0),
			handlers: make(map[string]HandlerFunc, 0)},
		validators: make(map[string]ValidatorFunc, 0)}
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Find the node for this path
	route, vars := ro.routes.resolveRoute(r.Method, r.URL.Path)

	// No registered route found
	if route == nil {
		http.NotFound(w, r)
		return
	}

	req := &Request{W: w, R: r}

	// Create a slice of variables to add to the Request object
	retVars := make(map[string]interface{})
	for i := 0; i < len(vars); i += 2 {
		// Check if there is a validator for this variable
		validator, ok := ro.validators[vars[i]]
		if ok {
			// Call the validator, if not succeed, send 404 code
			ok, variable := validator(vars[i+1])
			if !ok {
				http.NotFound(w, r)
				return
			}
			// Or add it to the slice
			retVars[vars[i]] = variable
		} else {
			// Add unvalidated variable to slice
			retVars[vars[i]] = vars[i+1]
		}
	}
	req.Vars = retVars

	// Check if there is a handler for the HTTP method
	handler, ok := route.handlers[r.Method]
	if ok {
		// Call the handler
		handler(req)
		return
	}

	// No handler registered
	http.NotFound(w, r)
}

// Handle registers a handler function for the specified HTTP method and URL pattern
func (ro *Router) Handle(method string, path string, handler HandlerFunc) {
	ro.routes.addRoute(method, path, handler)
}

// Validate registers a validation function for the specified dynamic variable
func (ro *Router) Validate(varName string, validator ValidatorFunc) {
	ro.validators[varName] = validator
}
