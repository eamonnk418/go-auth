package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eamonnk418/go-auth/internal/config"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

const (
	key    = "randomString" // make it more secure in production then this
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("config", cfg)

	store := sessions.NewCookieStore([]byte(key))

	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd
	store.Options.SameSite = http.SameSiteNoneMode

	gothic.Store = store

	goth.UseProviders(
		github.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL),
		// other providers can go here
	)

}
