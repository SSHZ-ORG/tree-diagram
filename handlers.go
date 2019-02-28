package main

import (
	"net/http"

	"github.com/SSHZ-ORG/tree-diagram/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func main() {
	r := mux.NewRouter()

	handlers.RegisterAPI(r)
	handlers.RegisterCommand(r)
	handlers.RegisterCrawl(r)
	handlers.RegisterCron(r)

	http.Handle("/", r)
	appengine.Main()
}
