package example1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	r.Use(basicLogging)

	r.Get("/", s.handleIndex)

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticFolder))))

	r.Group(func(r chi.Router) {
		//r.Use(authMiddleware1)
		r.Use(authMiddleware2([]userPass{{"admin", "admin"}, {"test1", "test1"}}))

		r.Get("/secret", s.handleSecret)
	})

	r.Post("/json", s.handleJSON)

	return r
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello b!"))
}

func (s *Server) handleSecret(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("This is a secret!"))
}

func (s *Server) handleJSON(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Apple string `json:"A"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	//b, _ := ioutil.ReadAll(r.Body)
	//err := json.Unmarshal(b, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(struct {
		A string `json:"a"`
		B string `json:"b"`
	}{
		A: input.Apple,
		B: "some value",
	})
	//b, _ := json.Marshal(struct {
	//	A string `json:"a"`
	//	B string `json:"b"`
	//}{
	//	A: input.Apple,
	//	B: "some value",
	//})
	//_, _ = w.Write(b)
}
