package example2

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server .
type Server struct {
	app AppInterface
}

// NewHTTP .
func NewHTTP(app AppInterface) *Server {
	return &Server{
		app: app,
	}
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", s.handleIndex)
	r.Get("/greeting", s.handleGreeting)

	return r
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	out, err := s.app.Index(context.Background(), indexInput{})
	if err != nil {
		_, _ = w.Write([]byte("Error happened"))
		return
	}

	_, _ = w.Write([]byte(out.message))
}

func (s *Server) handleGreeting(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	out, err := s.app.Greeting(context.Background(), greetingInput{name: q.Get("name")})
	if err != nil {
		_, _ = w.Write([]byte("Error happened"))
		return
	}

	_, _ = w.Write([]byte(out.message))
}
