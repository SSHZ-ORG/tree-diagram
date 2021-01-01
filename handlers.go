package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/SSHZ-ORG/tree-diagram/handlers"
	"github.com/SSHZ-ORG/tree-diagram/handlers/api"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := httprouter.New()
	handlers.RegisterCommand(r)
	handlers.RegisterCrawl(r)
	handlers.RegisterCron(r)

	grpc := api.GrpcServer()

	http.Handle("/", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if grpc.IsGrpcWebRequest(req) || grpc.IsAcceptableGrpcCorsRequest(req) {
			grpc.ServeHTTP(resp, req)
			return
		}
		r.ServeHTTP(resp, req)
	}))
	appengine.Main()
}
