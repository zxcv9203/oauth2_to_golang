package main

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	authServerURL = "http://localhost:8080"
)

var (
	config = oauth2.Config{
		ClientID:     "id",
		ClientSecret: "pw",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/token",
		},
	}
	globalToken *oauth2.Token
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, u, http.StatusFound)
	})
	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		code := r.Form.Get("code")
		token, err :=
			config.Exchange(context.Backgroun(),
				code, oauth2.SetAuthURLParam("code_verifier",
					"s256example"))
	})
}
