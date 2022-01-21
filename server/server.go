package main

import (
	"flag"
	"log"

	"github.com/go-oauth2/oauth2/store"
	"github.com/go-oauth2/oauth2/v4/manage"
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
	flag.StringVar(&domainvar, "r", "http://localhost:8080", "The domain of the redirect URL")
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
}
