// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorest

import (
	"net/http"
	"regexp"
)

type router struct {
	routes       []*Route
	defaultRoute http.Handler
}

func newRouter() *router {
	rts := []*Route{}
	return &router{routes: rts, defaultRoute: http.NotFoundHandler()}
}

func (rtr *router) handle(path string, handler http.Handler) *Route {
	methods := []string{}
	re := regexp.MustCompile(path)
	rt := &Route{re: re, methods: methods, handler: handler}
	rtr.routes = append(rtr.routes, rt)
	return rt
}

// ServeHTTP will execute the right handler depending on the route.
func (rtr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, rt := range rtr.routes {
		if rt.matchesPath(r.URL.Path) && rt.matchesMethod(r.Method) {
			rt.handler.ServeHTTP(w, r)
			return
		}
	}
	rtr.defaultRoute.ServeHTTP(w, r)
}

// A Route represents a REST route.
type Route struct {
	re      *regexp.Regexp
	methods []string
	handler http.Handler
}

// Methods allows to filter which HTTP methods will be handled by a given
// route. e.g.: "GET", "POST", "PUT", etc.
func (rt *Route) Methods(methods ...string) *Route {
	rt.methods = methods
	return rt
}

func (rt *Route) matchesPath(path string) bool {
	return rt.re.MatchString(path)
}

func (rt *Route) matchesMethod(method string) bool {
	if len(rt.methods) == 0 {
		return true
	}
	for _, m := range rt.methods {
		if m == method {
			return true
		}
	}
	return false
}
