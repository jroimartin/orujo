package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	restlog "github.com/jroimartin/gorest/handlers/log"
	"github.com/jroimartin/gorest/server"
)

func main() {
	s := server.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags)
	logHandler := server.M(restlog.NewLogHandler(logger,
		"{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}"))

	s.Route("/hello/{name}", server.H(http.HandlerFunc(helloHandler)), logHandler)
	s.Route("/test", server.H(http.HandlerFunc(testHandler)), logHandler)

	log.Fatalln(s.ListenAndServe())
}

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(rw, "Hello,", vars["name"])
}

func testHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintln(rw, "401 unauthorized")
}
