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

	p := NewPipe(
		http.HandlerFunc(h1),
		http.HandlerFunc(h2),
		http.HandlerFunc(h3),
	)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	p.ServeHTTP(rec, req)
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

	p := NewPipe(
		http.HandlerFunc(h1),
		http.HandlerFunc(h2),
		M(http.HandlerFunc(h3)),
	)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	p.ServeHTTP(rec, req)
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

	p := NewPipe(
		http.HandlerFunc(h1),
		http.HandlerFunc(h2),
	)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	p.ServeHTTP(rec, req)
	if len(want) != len(regErrors) {
		t.Fatalf("len(Errors(w))=%d; want=%d", len(regErrors), len(want))
	}
	for i := range regErrors {
		if regErrors[i].Error() != want[i].Error() {
			t.Errorf("Errors(w)[%d]=%s; want=%s", i, regErrors[i], want[i])
		}
	}
}
