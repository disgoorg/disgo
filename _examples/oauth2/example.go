package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/json"
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
	if err == nil {
		session := client.SessionController().GetSession(cookie.Value)
		if session != nil {
			var user *oauth2.User
			user, err = client.GetUser(session)
			if err != nil {
				writeError(w, "error while starting session", err)
				return
			}
			var userJSON []byte
			userJSON, err = json.MarshalIndent(user, "<br />", "&ensp;")
			if err != nil {
				writeError(w, "error while starting session", err)
				return
			}
			body = fmt.Sprintf("user:<br />%s", userJSON)
		}
	}
	if body == "" {
		body = `<button><a href="/login">login</a></button>`
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(body))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, client.GenerateAuthorizationURL(baseURL+"/trylogin", discord.ApplicationScopeIdentify, discord.ApplicationScopeGuilds, discord.ApplicationScopeEmail, discord.ApplicationScopeConnections), http.StatusMovedPermanently)
}

func handleTryLogin(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL.Query()
		code  = query.Get("code")
		state = query.Get("state")
	)
	if code != "" && state != "" {
		identifier := oauth2.RandStr(32)
		_, err := client.StartSession(code, state, identifier)
		if err != nil {
			writeError(w, "error while starting session", err)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "token", Value: identifier})
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func writeError(w http.ResponseWriter, text string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(text + ": " + err.Error()))
}
