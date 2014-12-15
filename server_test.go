// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorest

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	checks := []struct {
		path   string
		method string
		want   []byte
	}{
		{"/h1", "GET", []byte("h1")},
		{"/h1", "POST", []byte("h1")},
		{"/h1", "PUT", []byte("h3")},
		{"/h2", "GET", []byte("h2")},
		{"/h2", "POST", []byte("h2")},
		{"/h2", "PUT", []byte("h2")},
		{"/unk", "GET", []byte("h3")},
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
	s.Route("/h1", http.HandlerFunc(h1)).Methods("GET", "POST")
	s.Route("/h2", http.HandlerFunc(h2))
	s.RouteDefault(http.HandlerFunc(h3))
	ts := httptest.NewServer(s.mux)

	for _, check := range checks {
		urlStr := ts.URL + check.path
		client := &http.Client{}
		req, err := http.NewRequest(check.method, urlStr, nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.Equal(result, check.want) {
			t.Errorf("Server(%s, %s)=%s; want=%s",
				check.path, check.method, result, check.want)
		}
	}
}
