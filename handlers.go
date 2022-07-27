package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/SSHZ-ORG/tree-diagram/handlers"
	"github.com/SSHZ-ORG/tree-diagram/handlers/api"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := httprouter.New()
	handlers.RegisterCommand(r)
	handlers.RegisterCrawl(r)
	handlers.RegisterCron(r)

	grpc := api.GrpcServer()
	gzip := gziphandler.GzipHandler(grpc)

	http.Handle("/", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if grpc.IsAcceptableGrpcCorsRequest(req) {
			grpc.ServeHTTP(resp, req)
			return
		}
		if grpc.IsGrpcWebRequest(req) {
			gzip.ServeHTTP(resp, req)
			return
		}
		r.ServeHTTP(resp, req)
	}))
	appengine.Main()
}
