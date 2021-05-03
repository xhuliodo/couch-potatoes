package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/xhuliodo/couch-potatoes/clean-api/common/commands"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure"
)

func main() {
	log.Println("po ngrihet avioni...")
	ctx := commands.Context()

	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "letmein", "neo4j"))
	if err != nil {
		log.Fatal("db config is fucked up")
	}
	router := createApp(driver)

	// router := createApp()

	server := &http.Server{Addr: ":4000", Handler: router}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	log.Printf("App is listening on port: %s", server.Addr)

	<-ctx.Done()
	log.Println("App is shutting down")

	if err := server.Close(); err != nil {
		panic(err)
	}
}

func createApp(driver neo4j.Driver) (r *chi.Mux) {
	repo := infrastructure.NewNeo4jRepository(driver)

	r = infrastructure.CreateRouter()

	infrastructure.CreateRoutes(r, repo)

	return r
}

// func createApp() (r *chi.Mux) {
// 	repo := infrastructure.NewInMemoryRepository()

// 	r = commands.CreateRouter()

// 	interfaces.AddRoutes(r, repo, repo)

// 	return r
// }
