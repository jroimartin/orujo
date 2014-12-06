package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr   string
	mux    *http.ServeMux
	router *mux.Router
}

func NewServer(addr string) *Server {
	s := new(Server)
	s.Addr = addr
	s.mux = http.NewServeMux()
	s.router = mux.NewRouter()
	s.mux.Handle("/", s.router)
	return s
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.Addr, s.mux)
}

func (s *Server) ListenAndServeTLS(certFile string, keyFile string) error {
	return http.ListenAndServeTLS(s.Addr, certFile, keyFile, s.mux)
}

type Route mux.Route

func (s *Server) Route(path string, handlers ...Handler) *Route {
	return (*Route)(s.router.Handle(path, newPipe(handlers...)))
}

func (r *Route) Methods(methods ...string) *Route {
	return r.Methods(methods...)
}

type Handler interface {
	SetMandatory(v bool)
	Mandatory() bool
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
}

func M(h Handler) Handler {
	h.SetMandatory(true)
	return h
}
