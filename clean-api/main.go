package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
	"github.com/xhuliodo/couch-potatoes/clean-api/common/commands"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/auth"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/logger"
)

// @title Couch Potatoes clean API
// @version 1.0
// @description more movie more problems
// @host api.cp.dev.cloudapp.al
// host localhost:4001
// @basepath /
// @schemes https
// schemes http
// @accept json
// @produce json
// @contact.name Xhulio Doda
// @contact.url https://www.linkedin.com/in/xhulio-doda-745b41164/
// @contact.email xhuliodo@gmail.com

// gracefulShutdownTimeout bounds how long in-flight requests may finish after
// SIGINT/SIGTERM before the process exits.
const gracefulShutdownTimeout = 10 * time.Second

func main() {
	log.Println("po ngrihet avioni...")
	ctx := commands.Context()

	driver, err := neo4j.NewDriver(
		getenv("NEO4J_URI", "bolt://localhost:7687"),
		neo4j.BasicAuth(
			getenv("NEO4J_USER", "neo4j"),
			getenv("NEO4J_PASSWORD", "letmein"),
			"neo4j"), func(c *neo4j.Config) {
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
	server := &http.Server{Addr: getenv("API_LISTEN_PORT", ":4001"), Handler: router}
	serverErr := make(chan error, 1)

	if err := auth.CacheJwksCert(); err != nil {
		errorLogger.Log(err)
		log.Fatalf("jwks bootstrap failed: %v", err)
	}

	go func() {
		serverErr <- server.ListenAndServe()
	}()
	log.Printf("App is listening on port: %s", server.Addr)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("api server failed: %v", err)
		}
	case <-ctx.Done():
	}

	log.Println("App is shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		errorLogger.Log(err)
	}

	if err := driver.Close(); err != nil {
		errorLogger.Log(err)
	}

	log.Println("App stopped")
}

func createApp(driver neo4j.Driver, accessLogger *logrus.Logger, errorLogger *logger.ErrorLogger) (r *chi.Mux) {
	repo := db.NewNeo4jRepository(driver)

	r = infrastructure.CreateRouter(accessLogger)

	infrastructure.CreateRoutes(r, repo, errorLogger)

	return r
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
