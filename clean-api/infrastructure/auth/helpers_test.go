package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

const testJwks = `{
  "keys": [
    {
      "kty": "RSA",
      "kid": "key-1",
      "use": "sig",
      "n": "n-value",
      "e": "AQAB",
      "x5c": ["cert-one"]
    }
  ]
}`

// newJwksServer serves a fixed JWKS body with the given status.
func newJwksServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func cachePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "auth", "jwks.json")
}

func TestCacheJwksCertWritesCacheWithRestrictedPermissions(t *testing.T) {
	path := cachePath(t)
	t.Setenv("JWKS_URL", newJwksServer(t, http.StatusOK, testJwks).URL)
	t.Setenv("JWKS_CACHE_PATH", path)

	if err := CacheJwksCert(); err != nil {
		t.Fatalf("CacheJwksCert() error = %v", err)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("cache file not written: %v", err)
	}
	if got := fileInfo.Mode().Perm(); got != 0o600 {
		t.Errorf("cache file mode = %04o, want 0600", got)
	}

	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("cache dir missing: %v", err)
	}
	if got := dirInfo.Mode().Perm(); got != 0o700 {
		t.Errorf("cache dir mode = %04o, want 0700", got)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read cache: %v", err)
	}
	var cached Jwks
	if err := json.Unmarshal(data, &cached); err != nil {
		t.Fatalf("cached file is not valid JWKS JSON: %v", err)
	}
	if len(cached.Keys) != 1 || cached.Keys[0].Kid != "key-1" {
		t.Errorf("cached keys = %+v, want key-1", cached.Keys)
	}
}

func TestCacheJwksCertTightensExistingPermissions(t *testing.T) {
	path := cachePath(t)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("prepare dir: %v", err)
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatalf("prepare dir mode: %v", err)
	}
	if err := os.WriteFile(path, []byte(`{"keys":[]}`), 0o644); err != nil {
		t.Fatalf("prepare stale cache: %v", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatalf("prepare stale cache mode: %v", err)
	}

	t.Setenv("JWKS_URL", newJwksServer(t, http.StatusOK, testJwks).URL)
	t.Setenv("JWKS_CACHE_PATH", path)

	if err := CacheJwksCert(); err != nil {
		t.Fatalf("CacheJwksCert() error = %v", err)
	}

	if info, err := os.Stat(dir); err != nil {
		t.Fatalf("stat dir: %v", err)
	} else if got := info.Mode().Perm(); got != 0o700 {
		t.Errorf("cache dir mode = %04o, want 0700", got)
	}
	if info, err := os.Stat(path); err != nil {
		t.Fatalf("stat file: %v", err)
	} else if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("cache file mode = %04o, want 0600", got)
	}
}

func TestCacheJwksCertRejectsBadResponses(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
	}{
		{name: "non 200 status", status: http.StatusInternalServerError, body: `{"keys":[{"kid":"key-1","x5c":["c"]}]}`},
		{name: "not json", status: http.StatusOK, body: `<html>nope</html>`},
		{name: "no keys", status: http.StatusOK, body: `{"keys":[]}`},
		{name: "key without x5c", status: http.StatusOK, body: `{"keys":[{"kid":"key-1"}]}`},
		{name: "key without kid", status: http.StatusOK, body: `{"keys":[{"x5c":["cert"]}]}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := cachePath(t)
			t.Setenv("JWKS_URL", newJwksServer(t, tc.status, tc.body).URL)
			t.Setenv("JWKS_CACHE_PATH", path)

			if err := CacheJwksCert(); err == nil {
				t.Fatal("CacheJwksCert() error = nil, want non-nil")
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Errorf("cache file written on failure: stat err = %v", err)
			}
		})
	}
}

func TestCacheJwksCertRejectsUnreachableURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close() // nothing is listening anymore

	path := cachePath(t)
	t.Setenv("JWKS_URL", url)
	t.Setenv("JWKS_CACHE_PATH", path)

	if err := CacheJwksCert(); err == nil {
		t.Fatal("CacheJwksCert() error = nil, want non-nil")
	}
}

func TestGetPemCertReturnsErrorWhenCacheUnreadable(t *testing.T) {
	t.Setenv("JWKS_CACHE_PATH", filepath.Join(t.TempDir(), "missing", "jwks.json"))

	cert, err := GetPemCert(&jwt.Token{Header: map[string]interface{}{"kid": "key-1"}})
	if err == nil {
		t.Fatalf("GetPemCert() = %q, nil; want error", cert)
	}
	if !strings.Contains(err.Error(), "jwks") {
		t.Errorf("error %q should mention the jwks cache", err)
	}
}

func TestGetPemCertReturnsPEMForKnownKid(t *testing.T) {
	path := cachePath(t)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(testJwks), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JWKS_CACHE_PATH", path)

	cert, err := GetPemCert(&jwt.Token{Header: map[string]interface{}{"kid": "key-1"}})
	if err != nil {
		t.Fatalf("GetPemCert() error = %v", err)
	}
	if !strings.Contains(cert, "cert-one") {
		t.Errorf("cert = %q, want it to contain cert-one", cert)
	}

	if _, err := GetPemCert(&jwt.Token{Header: map[string]interface{}{"kid": "unknown"}}); err == nil {
		t.Error("GetPemCert() error = nil for unknown kid, want non-nil")
	}
}

// signToken builds a syntactically valid JWT with the given kid. The signature
// is intentionally not verifiable by the JWKS keys; these tests exercise the
// path up to key lookup, not signature verification.
func signToken(t *testing.T, kid string) string {
	t.Helper()
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["kid"] = kid
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

func TestAuthMiddlewareUnauthorizedWhenJwksCacheUnreadable(t *testing.T) {
	t.Setenv("JWKS_CACHE_PATH", filepath.Join(t.TempDir(), "missing", "jwks.json"))

	handlerCalled := false
	handler := AuthMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("authorization", "Bearer "+signToken(t, "key-1"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if handlerCalled {
		t.Error("protected handler was called without a readable JWKS cache")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
