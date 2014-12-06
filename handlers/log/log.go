package log

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/jroimartin/gorest/server"
)

type LogHandler struct {
	log       *log.Logger
	tmpl      *template.Template
	mandatory bool
}

func NewLogHandler(logger *log.Logger, fmt string) server.Handler {
	tmpl := template.Must(template.New("fmt").Parse(fmt))
	return &LogHandler{log: logger, tmpl: tmpl}
}

func (h *LogHandler) SetMandatory(v bool) {
	h.mandatory = v
}

func (h *LogHandler) Mandatory() bool {
	return h.mandatory
}

func (h *LogHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	data := struct {
		Resp http.ResponseWriter
		Req  *http.Request
	}{rw, r}

	var out bytes.Buffer
	h.tmpl.Execute(&out, data)
	h.log.Println(out.String())
}
