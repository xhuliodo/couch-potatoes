package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const testKid = "key-1"

// writeTestJwks generates a signing key and a JWKS cache holding its
// self-signed certificate, mirroring the Auth0 response shape.
func writeTestJwks(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "couch-potatoes-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	jwks := Jwks{Keys: []JSONWebKeys{{
		Kty: "RSA",
		Kid: testKid,
		Use: "sig",
		X5c: []string{base64.StdEncoding.EncodeToString(der)},
	}}}
	data, err := json.Marshal(jwks)
	if err != nil {
		t.Fatalf("marshal jwks: %v", err)
	}

	path := filepath.Join(t.TempDir(), "jwks.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write jwks: %v", err)
	}
	t.Setenv("JWKS_CACHE_PATH", path)

	return key
}

func signRS256(t *testing.T, key *rsa.PrivateKey, kid string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	key := writeTestJwks(t)

	var gotUserID interface{}
	var gotIsAdmin interface{}
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID = r.Context().Value("userId")
		gotIsAdmin = r.Context().Value("isAdmin")
		w.WriteHeader(http.StatusOK)
	}))

	token := signRS256(t, key, testKid, jwt.MapClaims{
		"sub":                                "auth0|user-1",
		"https://couch-potatoes.com/claims/": []interface{}{"admin"},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body %s)", rec.Code, http.StatusOK, rec.Body)
	}
	if gotUserID != "auth0|user-1" {
		t.Errorf("userId = %v, want auth0|user-1", gotUserID)
	}
	if gotIsAdmin != true {
		t.Errorf("isAdmin = %v, want true", gotIsAdmin)
	}
}

func TestAuthMiddlewareRejectsInvalidTokens(t *testing.T) {
	key := writeTestJwks(t)
	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate other key: %v", err)
	}

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cases := []struct {
		name  string
		token string
	}{
		{
			name:  "signed by another key",
			token: signRS256(t, otherKey, testKid, jwt.MapClaims{"sub": "auth0|user-1"}),
		},
		{
			name:  "unknown kid",
			token: signRS256(t, key, "unknown-key", jwt.MapClaims{"sub": "auth0|user-1"}),
		},
		{
			name: "expired",
			token: signRS256(t, key, testKid, jwt.MapClaims{
				"sub": "auth0|user-1",
				"exp": time.Now().Add(-time.Hour).Unix(),
			}),
		},
		{
			name:  "not a token",
			token: "garbage",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("authorization", "Bearer "+tc.token)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
		})
	}
}
