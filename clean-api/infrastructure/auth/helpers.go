package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	defaultJwksURL  = "https://dev-ps5dqqis.eu.auth0.com/.well-known/jwks.json"
	defaultJwksPath = "/auth/jwks.json"
	jwksHTTPTimeout = 10 * time.Second
	maxJwksBytes    = 1 << 20
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func GetPemCert(token *jwt.Token) (
	cert string, err error,
) {
	path := jwksCachePath()

	data, err := os.ReadFile(path)
	if err != nil {
		return cert, fmt.Errorf("read jwks cache %s: %w", path, err)
	}

	var jwks = Jwks{}
	err = json.Unmarshal(data, &jwks)
	if err != nil {
		return cert, fmt.Errorf("decode jwks cache %s: %w", path, err)
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid && len(jwks.Keys[k].X5c) > 0 {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, fmt.Errorf("unable to find appropriate key for kid %v", token.Header["kid"])
	}

	return cert, nil
}

// CacheJwksCert fetches the tenant's JWKS and persists it to the cache path.
// It returns an error instead of exiting so the caller decides how to report
// and terminate.
func CacheJwksCert() error {
	jwksURL := getenv("JWKS_URL", defaultJwksURL)
	if _, err := url.ParseRequestURI(jwksURL); err != nil {
		return fmt.Errorf("invalid JWKS_URL %q: %w", jwksURL, err)
	}

	client := &http.Client{Timeout: jwksHTTPTimeout}
	resp, err := client.Get(jwksURL)
	if err != nil {
		return fmt.Errorf("fetch JWKS %s: %w", jwksURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch JWKS %s: unexpected HTTP status %d", jwksURL, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxJwksBytes))
	if err != nil {
		return fmt.Errorf("read JWKS %s: %w", jwksURL, err)
	}

	jwks := Jwks{}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return fmt.Errorf("decode JWKS %s: %w", jwksURL, err)
	}
	if len(jwks.Keys) == 0 {
		return fmt.Errorf("JWKS %s: response contains no keys", jwksURL)
	}
	for _, key := range jwks.Keys {
		if key.Kid == "" || len(key.X5c) == 0 {
			return fmt.Errorf("JWKS %s: key is missing kid or x5c", jwksURL)
		}
	}

	path := jwksCachePath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create jwks cache dir %s: %w", dir, err)
	}
	// MkdirAll does not tighten an existing directory.
	if err := os.Chmod(dir, 0o700); err != nil {
		return fmt.Errorf("set jwks cache dir permissions on %s: %w", dir, err)
	}

	if err := writeFileAtomic(path, body, 0o600); err != nil {
		return fmt.Errorf("write jwks cache %s: %w", path, err)
	}

	return nil
}

// writeFileAtomic writes data to path with perm, replacing any existing file
// only once the full content is safely on disk.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, path)
}

func jwksCachePath() string {
	return getenv("JWKS_CACHE_PATH", defaultJwksPath)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
