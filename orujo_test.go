// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orujo

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	checks := []struct {
		path   string
		method string
		want   string
	}{
		{"/h1", "GET", "h1"},
		{"/h1", "POST", "h1"},
		{"/h1", "PUT", "h3"},
		{"/h1x", "GET", "h3"},
		{"/h2", "GET", "h2"},
		{"/h2", "POST", "h2"},
		{"/h2", "PUT", "h2"},
		{"/h2x", "GET", "h2"},
		{"/unk", "GET", "h3"},
	}

	h1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("h1"))
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("h2"))
	}
	h3 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("h3"))
	}

	s := NewServer("")
	s.Route(`^/h1$`, http.HandlerFunc(h1)).Methods("GET", "POST")
	s.Route(`^/h2`, http.HandlerFunc(h2))
	s.RouteDefault(http.HandlerFunc(h3))

	for _, check := range checks {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(check.method, check.path, nil)
		if err != nil {
			log.Fatal(err)
		}
		s.mux.ServeHTTP(rec, req)
		result := rec.Body.String()
		if result != check.want {
			t.Errorf("Router(%s, %s)=%s; want=%s",
				check.method, check.path, result, check.want)
		}
	}
}

func TestPipeQuit(t *testing.T) {
	want := "h1h2"
	result := ""

	h1 := func(w http.ResponseWriter, r *http.Request) {
		result += "h1"
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		result += "h2"
		w.WriteHeader(401)
	}
	h3 := func(w http.ResponseWriter, r *http.Request) {
		result += "h3"
	}

	s := NewServer("")
	s.Route(`.*`,
		http.HandlerFunc(h1),
		http.HandlerFunc(h2),
		http.HandlerFunc(h3),
	)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	s.mux.ServeHTTP(rec, req)
	if result != want {
		t.Errorf("Pipe(h1, h2, h3)=%s; want=%s", result, want)
	}
}

func TestPipeMandatory(t *testing.T) {
	want := "h1h2h3"
	result := ""

	h1 := func(w http.ResponseWriter, r *http.Request) {
		result += "h1"
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		result += "h2"
		w.WriteHeader(401)
	}
	h3 := func(w http.ResponseWriter, r *http.Request) {
		result += "h3"
	}

	s := NewServer("")
	s.Route(`.*`,
		http.HandlerFunc(h1),
		http.HandlerFunc(h2),
		M(http.HandlerFunc(h3)),
	)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	s.mux.ServeHTTP(rec, req)
	if result != want {
		t.Errorf("Pipe(h1, h2, M(h3))=%s; want=%s", result, want)
	}
}

func TestErrors(t *testing.T) {
	want := []error{
		errors.New("Err1.1"),
		errors.New("Err1.2"),
		errors.New("Err2.1"),
		errors.New("Err2.3"),
	}

	var regErrors []error
	h1 := func(w http.ResponseWriter, r *http.Request) {
		RegisterError(w, errors.New("Err1.1"))
		RegisterError(w, errors.New("Err1.2"))
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		RegisterError(w, errors.New("Err2.1"))
		RegisterError(w, errors.New("Err2.3"))
		regErrors = Errors(w)
	}

	s := NewServer("")
	s.Route(`.*`,
		http.HandlerFunc(h1),
		http.HandlerFunc(h2),
	)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	s.mux.ServeHTTP(rec, req)
	if len(want) != len(regErrors) {
		t.Fatalf("len(Errors(w))=%d; want=%d", len(regErrors), len(want))
	}
	for i := range regErrors {
		if regErrors[i].Error() != want[i].Error() {
			t.Errorf("Errors(w)[%d]=%s; want=%s", i, regErrors[i], want[i])
		}
	}
}
