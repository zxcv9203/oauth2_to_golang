package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-oauth2/oauth2/models"
	"github.com/go-oauth2/oauth2/store"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/errors"
)

var (
	dumpvar   bool
	idvar     string
	secretvar string
	domainvar string
	portvar   int
)

func init() {
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.StringVar(&idvar, "i", "OAuthID-example", "The Client id being passed in")
	flag.StringVar(&secretvar, "s", "OAuthSecret-example", "The Client Secret being passed in")
	flag.StringVar(&domainvar, "req", "http://localhost:8080", "The domain of the redirect URL")
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {
	flag.Parse()
	if dumpvar {
		log.Println("Dumping requests")
	}
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()
	clientStore.Set(idvar, &models.Client)

	srv := server.NewServer(server.NewConfg(), manager)

	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (req *errors.Response) {
		log.Println("Internal Error : ", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(req *errors.Response) {
		log.Println("Response Error : ", req.Error.Error())
		return
	})

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)

	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, req *http.Request) {
		if dumpvar {
			dumpRequest(os.Stdout, "authorize", req)
		}
		store, err := session.Start(req.Context(), w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		req.Form = form
		store.Delete("ReturnUri")
		store.Save()

		err = srv.HandleAuthorizeRequest(w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
	http.HandleFunc("/oauth/token", func(w http.ResponseWriter, req *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "token", req)
		}
		err := srv.HandleTokenRequest(w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "test", req)
		}
		token, err := srv.ValidationBearerToken(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpireIn()).Sub(time.Now()).Seconds),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", " ")
		e.Encode(data)
	})
	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil))
}

func dumpRequest(w io.Writer, header string, req *http.Request) error {
	data, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	w.Write([]byte("\n" + header + " : \n"))
	w.Write(data)
	return nil
}

func userAuthorizeHandler(w http.ResponseWriter, req *http.Request) (userID string, err error) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", req)
	}
	store, err := session.Start(req.Context(), w, req)
	if err != nil {
		return
	}
	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if req.Form == nil {
			req.ParseForm()
		}
		store.Set("ReturnUri", req.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "login", req)
	}
	store, err := session.Start(req.Context(), w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.Method == "post" {
		if req.Form == nil {
			if err != req.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", req.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, req, "static/login.html")
}

func authHandler(w http.ResponseWriter, req *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "auth", req)
	}
	store, err := session.Start(nil, w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, req, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
