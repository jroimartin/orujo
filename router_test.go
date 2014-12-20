// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orujo

import (
	"net/http"
	"testing"
)

func TestVars(t *testing.T) {
	expr := `/name/(?P<name>\w+)/id/(?P<id>\d+)`
	want := make(map[string]string)
	want["name"] = "go"
	want["id"] = "60"

	r, err := http.NewRequest("GET", "http://localhost/name/go/id/60", nil)
	if err != nil {
		t.Fatal(err)
	}

	result := Vars(r, expr)

	if result["name"] != want["name"] {
		t.Errorf("vars[\"name\"]=%s; want=%s",
			result["name"], want["name"])
	}
	if result["id"] != want["id"] {
		t.Errorf("vars[\"id\"]=%s; want=%s",
			result["id"], want["id"])
	}
}
