// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasic(t *testing.T) {
	logLine := `{{.Req.Method}} {{.Req.URL.Path}}`
	want := "[LOG] GET /h1\n"

	logBuffer := new(bytes.Buffer)
	logger := stdlog.New(logBuffer, "[LOG] ", 0)
	h := NewLogHandler(logger, logLine)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost/h1", nil)
	if err != nil {
		stdlog.Fatal(err)
	}
	h.ServeHTTP(rec, req)

	result := logBuffer.String()
	if result != want {
		t.Errorf("Log=%s; want=%s", result, want)
	}
}
