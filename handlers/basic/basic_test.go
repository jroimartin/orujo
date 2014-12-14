// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasic(t *testing.T) {
	checks := []struct {
		username string
		password string
		want     int
	}{
		{"non-existent", "non-existent", http.StatusUnauthorized},
		{"username", "non-existent", http.StatusUnauthorized},
		{"non-existent", "password123", http.StatusUnauthorized},
		{"user", "password123", http.StatusOK},
	}
	h := NewBasicHandler("basic", "user", "password123")
	for _, check := range checks {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(check.username, check.password)
		h.ServeHTTP(rec, req)
		if rec.Code != check.want {
			t.Errorf("Basic(%s, %s)=%d; want=%d",
				check.username, check.password, rec.Code, check.want)
		}
	}
}
