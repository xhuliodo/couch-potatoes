package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
	"github.com/xhuliodo/couch-potatoes/clean-api/common/commands"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/logger"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
)

// @title Couch Potatoes clean API
// @version 1.0
// @description more movie more problems
// host api.cp.dev.cloudapp.al
// @host localhost:4000
// @basepath /
// schemes https
// @schemes http
// @accept json
// @produce json
// @contact.name Xhulio Doda
// @contact.url https://www.linkedin.com/in/xhulio-doda-745b41164/
// @contact.email xhuliodo@gmail.com

func main() {
	log.Println("po ngrihet avioni...")
	ctx := commands.Context()

	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "letmein", "neo4j"), func(c *neo4j.Config) {
		time := time.Second * 5
		c.ConnectionAcquisitionTimeout = time
		c.MaxTransactionRetryTime = time * 2
	})
	if err != nil {
		log.Fatal("db config is fucked up")
	}

	accessLogger := logger.NewAccessLogger()
	errorLogger := logger.NewErrorLogger()

	router := createApp(driver, accessLogger, errorLogger)
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

func createApp(driver neo4j.Driver, accessLogger *logrus.Logger, errorLogger *logger.ErrorLogger) (r *chi.Mux) {
	repo := db.NewNeo4jRepository(driver)

	r = infrastructure.CreateRouter(accessLogger)

	infrastructure.CreateRoutes(r, repo, errorLogger)

	return r
}
