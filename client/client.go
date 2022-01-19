package main

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

/*
	oauth2란?
	구글, 페이스북, 카카오 등에서 제공하는 인증 서버를 통해
	회원 정보를 인증하고 Access Token을 발급받기 위한 표준 프로토콜

	발급받은 Access Token을 이용하여 인증 받은 곳(구글, 페이스북, 카카오 등)
	의 API 서비스를 이용할 수 있게 됩니다.
*/

/*
	oauth2 용어

	Resource owner : Resource server로 부터 계정을 소유하고 있는 사용자를 의미합니다.

	Client : 구글, 페이스북, 카카오 등의 API 서비스를 이용하는 제 3의 서비스를 의미합니다.

	Authorization Server(권한 서버) : 권한을 관리해주는 서버, Access
*/
const (
	authServerURL = "http://localhost:8080"
)

/*
	config 각 변수 역할
	ClientID : OAuth client
*/
var (
	config = oauth2.Config{
		ClientID:     "OAuthID-example",
		ClientSecret: "OAuthSecret-example",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := config.AuthcodeURL("xyz",
			oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		if state != "xyz" {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}
		code := r.Form.Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}
		token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		globalToken = token

		e := json.NewEncoder(w)
		e.SetIndent
	})
}
