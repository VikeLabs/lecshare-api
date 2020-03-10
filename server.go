package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vikelabs/lecshare-api/graph"
	"github.com/vikelabs/lecshare-api/graph/generated"
)

var (
	port int
	host string
)

func main() {
	flag.IntVar(&port, "p", 8080, "specify port to use")
	flag.StringVar(&host, "h", "localhost", "specfiy host to bind to")

	flag.Parse()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	// define routes for development.
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graph.Middleware(srv))

	// host on user defined port / default port.
	log.Printf("connect to http://%s:%d/ for GraphQL playground", host, port)
	log.Printf("connect to http://%s:%d/query for GraphQL endpoint", host, port)
	log.Fatal(http.ListenAndServe(host+":"+strconv.Itoa(port), nil))
}
