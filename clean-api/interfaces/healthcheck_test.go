package interfaces_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
	"github.com/xhuliodo/couch-potatoes/clean-api/interfaces"
)

type silentErrorLogger struct{}

func (silentErrorLogger) Log(error) {}

// unreachableDatabase keeps the in-memory fixtures but reports a database
// connection failure, the same cause the Neo4j repository produces.
type unreachableDatabase struct {
	*db.InMemoryRepository
}

func (unreachableDatabase) Healthcheck() error {
	cause := errors.New("db_connection")
	return errors.Wrap(cause, "connection refused")
}

func TestReadyReportsFailureWhenDatabaseIsUnreachable(t *testing.T) {
	repo := unreachableDatabase{db.NewInMemoryRepository()}
	healthcheck := interfaces.NewHealthcheckInterface(repo, silentErrorLogger{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	healthcheck.GetReady(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("GET /ready status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestReadyReportsSuccessWhenDatabaseIsReachable(t *testing.T) {
	repo := db.NewInMemoryRepository()
	healthcheck := interfaces.NewHealthcheckInterface(repo, silentErrorLogger{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	healthcheck.GetReady(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /ready status = %d, want %d", rec.Code, http.StatusOK)
	}
}
