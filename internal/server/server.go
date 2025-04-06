package server

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/eamonnk418/go-auth/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// Server holds the dependencies for our HTTP server.
type Server struct {
	DB     database.Database
	Router http.Handler
}

// NewServer creates a new Server instance with its dependencies injected.
func NewServer(db database.Database) *Server {
	s := &Server{
		DB: db,
	}
	// Register routes.
	s.Router = s.RegisterRoutes()
	return s
}

// Start runs the HTTP server on the given address.
func (s *Server) Start(addr string) error {
	log.Printf("Server running on %s", addr)
	return http.ListenAndServe(addr, s.Router)
}

// LoginHandler renders a simple login page with a GitHub login link.
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("login").Parse(loginTemplate)
	if err != nil {
		log.Printf("error parsing login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

// Callback handles the OAuth callback.
// It completes authentication, saves the user in session, and then redirects to the profile page.
func (s *Server) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	// Ensure provider is available in the context for Gothic.
	r = r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("error during CompleteUserAuth: %v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Save user details to session.
	session, err := gothic.Store.Get(r, gothic.SessionName)
	if err != nil {
		log.Printf("error getting session: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	session.Values["user"] = user
	if err := session.Save(r, w); err != nil {
		log.Printf("error saving session: %v", err)
		http.Error(w, "Session save error", http.StatusInternalServerError)
		return
	}

	// Redirect to the profile UI.
	http.Redirect(w, r, "/profile", http.StatusMovedPermanently)
}

// Profile renders the UI page showing the user's details.
// It reads the user info from the session.
func (s *Server) Profile(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, gothic.SessionName)
	if err != nil {
		log.Printf("error getting session: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	u, ok := session.Values["user"]
	if !ok {
		// User not logged in; redirect to login.
		http.Redirect(w, r, "/auth/github", http.StatusTemporaryRedirect)
		return
	}
	user, ok := u.(goth.User)
	if !ok {
		log.Printf("error asserting user type")
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	t, err := template.New("user").Parse(userTemplate)
	if err != nil {
		log.Printf("error parsing user template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, user)
}

// SignIn begins the OAuth flow.
func (s *Server) SignIn(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	// Set provider in context.
	r = r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider))
	gothic.BeginAuthHandler(w, r)
}

// SignOut logs the user out and redirects to the home/login page.
func (s *Server) SignOut(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Templates used for login and displaying user info.
var loginTemplate = `
<html>
<head>
    <title>Login</title>
</head>
<body>
    <h1>Login with GitHub</h1>
    <p>
      <a href="/auth/github">Login with GitHub</a>
    </p>
</body>
</html>
`

var userTemplate = `
<html>
<head>
    <title>User Info</title>
</head>
<body>
    <p><a href="/logout/{{.Provider}}">Logout</a></p>
    <h2>User Details</h2>
    <p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
    <p>Email: {{.Email}}</p>
    <p>NickName: {{.NickName}}</p>
    <p>Location: {{.Location}}</p>
    <p>Avatar: <img src="{{.AvatarURL}}" alt="avatar" width="50"></p>
    <p>Description: {{.Description}}</p>
    <p>UserID: {{.UserID}}</p>
    <p>AccessToken: {{.AccessToken}}</p>
    <p>ExpiresAt: {{.ExpiresAt}}</p>
    <p>RefreshToken: {{.RefreshToken}}</p>
</body>
</html>
`
