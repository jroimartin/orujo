// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sessions

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	h := NewSessionHandler("test", []byte("secret"))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	h.ServeHTTP(rec, req)

	sessionID, err := h.SessionID(req)
	if err != nil {
		log.Fatal(err)
	}

	if sessionID == "" {
		t.Errorf("SessionID=%s; want=<random ID>", sessionID)
	}
}
