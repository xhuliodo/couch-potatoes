package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/xhuliodo/couch-potatoes/clean-api/common/commands"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure"
	"github.com/xhuliodo/couch-potatoes/clean-api/interfaces"
)

func main() {
	log.Println("po ngrihet avioni...")
	ctx := commands.Context()

	router := createApp()

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

func createApp() (r *chi.Mux) {
	movieRepo := infrastructure.NewInMemoryRepository()

	r = commands.CreateRouter()

	interfaces.AddRoutes(r, movieRepo)

	return r
}
