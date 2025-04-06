package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Home page: show login button
	r.Get("/", s.LoginHandler)

	// Health endpoint
	r.Get("/health", s.healthHandler)

	// OAuth endpoints: begin auth, callback, and logout
	r.Get("/auth/{provider}", s.SignIn)
	r.Get("/auth/{provider}/callback", s.Callback)
	r.Get("/logout/{provider}", s.SignOut)

	// Profile UI page after login
	r.Get("/profile", s.Profile)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello, World!"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.DB.Health())
	_, _ = w.Write(jsonResp)
}
