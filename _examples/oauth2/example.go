package main

import (
	"net/http"
	"os"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/oauth2"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

var (
	clientID     = discord.Snowflake(os.Getenv("client_id"))
	clientSecret = os.Getenv("client_secret")
	logger       = log.Default()
	httpClient   = http.DefaultClient
	client       oauth2.Client

	sessions = map[string]oauth2.Session{}
)

func main() {
	logger.SetLevel(log.LevelDebug)
	logger.Info("starting ExampleBot...")
	logger.Infof("disgo %s", info.Version)

	client = oauth2.New(clientID, clientSecret, oauth2.WithLogger(logger), oauth2.WithRestClientConfigOpts(rest.WithHTTPClient(httpClient)))

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/login", handleLogin)
	http.ListenAndServe(":6969", mux)
}

func handleRoot(w http.ResponseWriter, request *http.Request) {
	var body string
	cookie, err := request.Cookie("token")
	if err != nil {
		body = `<button href="/login">login</button>`
	} else {

	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(body))
}

func handleLogin(w http.ResponseWriter, request *http.Request) {

}
