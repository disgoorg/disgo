package Disgo

import "net/http"

type Disgo struct {
	Token string
	Intents Intent
	Client *http.Client
}
