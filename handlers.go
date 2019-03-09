package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/SSHZ-ORG/tree-diagram/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()

	handlers.RegisterAPI(r)
	handlers.RegisterCommand(r)
	handlers.RegisterCrawl(r)
	handlers.RegisterCron(r)

	http.Handle("/", r)
	appengine.Main()
}
