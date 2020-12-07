package example1

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Server .
type Server struct {
	staticFolder string
}

// New .
func New(staticFolder string) *Server {
	return &Server{
		staticFolder: staticFolder,
	}
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	// TODO: add a basic logging middleware

	r.Get("/", s.handleIndex)

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticFolder))))

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware1)
		//r.Use(authMiddleware2([]userPass{{"admin", "admin"}, {"test1", "test1"}}))

		r.Get("/secret", s.handleSecret)
	})
	return r
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello b!"))
}

func (s *Server) handleSecret(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("This is a secret!"))
}
