package infrastructure_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure"
)

func preflight(t *testing.T, origin string) *httptest.ResponseRecorder {
	t.Helper()

	router := infrastructure.CreateRouter(discardLogger())
	router.Get("/probe", func(http.ResponseWriter, *http.Request) {})
	router.Options("/probe", func(http.ResponseWriter, *http.Request) {})

	req := httptest.NewRequest(http.MethodOptions, "/probe", nil)
	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func discardLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logger
}

func TestCORSAllowsOnlyOriginsFromTheEnvironment(t *testing.T) {
	t.Setenv("API_ALLOWED_ORIGINS", "https://couch-potatoes.test, https://admin.couch-potatoes.test")

	allowed := preflight(t, "https://admin.couch-potatoes.test")
	if got := allowed.Header().Get("Access-Control-Allow-Origin"); got != "https://admin.couch-potatoes.test" {
		t.Errorf("allowed origin: Access-Control-Allow-Origin = %q, want %q", got, "https://admin.couch-potatoes.test")
	}
	if got := allowed.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("allowed origin: Access-Control-Allow-Credentials = %q, want %q", got, "true")
	}
	if got := allowed.Header().Get("Access-Control-Allow-Methods"); got != http.MethodGet {
		t.Errorf("allowed origin: Access-Control-Allow-Methods = %q, want %q", got, http.MethodGet)
	}

	disallowed := preflight(t, "https://elsewhere.test")
	if got := disallowed.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("disallowed origin: Access-Control-Allow-Origin = %q, want no header", got)
	}
}

func TestCORSDefaultsToTheLocalClientOrigin(t *testing.T) {
	t.Setenv("API_ALLOWED_ORIGINS", "")

	rec := preflight(t, "http://localhost:3000")
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:3000")
	}
}
