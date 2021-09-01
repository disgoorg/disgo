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
	baseURL      = os.Getenv("base_url")
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
	mux.HandleFunc("/trylogin", handleTryLogin)
	http.ListenAndServe(":6969", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	var body string
	cookie, err := r.Cookie("token")
	if err != nil {
		body = `<button href="/login">login</button>`
	} else {

	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(body))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, client.GenerateAuthorizationURL(baseURL+"/trylogin"), http.StatusMovedPermanently)
}

func handleTryLogin(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL.Query()
		code  = query.Get("code")
		state = query.Get("state")
	)
	if code != "" && state != "" {
		session, err := client.StartSession(code, state, oauth2.RandStr(32))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("error while getting session: " + err.Error()))
			return
		}
		client.
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
