package server

import (
	"context"
	"log"
	"net/http"
	"regexp"
)

// Router holds the data about our Router.
type Router struct {
	Routes []Route
}

// Route holds the info about a route
type Route struct {
	Method  string
	Path    *regexp.Regexp
	Handler http.HandlerFunc
}

// Params takes in http.Request with its context,
// parses the path via regexp, and returns the parameter.
func Params(r *http.Request, name string) string {
	ctx := r.Context()

	params := ctx.Value("params").(map[string]string)

	return params[name]
}

// Route creates a new Route
func (rtr *Router) Route(method, path string, handler http.HandlerFunc) {
	rtr.Routes = append(rtr.Routes, Route{
		Method:  method,
		Path:    regexp.MustCompile(path),
		Handler: handler,
	})
}

// Match handles path validation.
// If no match, return false, otherwise, return true.
func (rt *Route) Match(r *http.Request) map[string]string {
	match := rt.Path.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return nil
	}

	params := make(map[string]string)
	gn := rt.Path.SubexpNames()
	for i, g := range match {
		params[gn[i]] = g
	}

	return params
}

// ServeHTTP serves the HTTP if everything is okay,
// else http.NotFound
func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("error:", r)
			http.Error(w, "something happened", http.StatusInternalServerError)
		}
	}()
	for _, rt := range rtr.Routes {
		params := rt.Match(r)
		if params == nil {
			continue
		}

		ctx := context.WithValue(r.Context(), "params", params)
		rt.Handler.ServeHTTP(w, r.WithContext(ctx))
		return
	}
	http.NotFound(w, r)
}
