package auth

import (
	"encoding/json"
	"io/ioutil"

	"github.com/form3tech-oss/jwt-go"
	"github.com/pkg/errors"
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
	// resp, err := http.Get("https://dev-ps5dqqis.eu.auth0.com/.well-known/jwks.json")
	// if err != nil {
	// 	return cert, err
	// }
	// defer resp.Body.Close()
	data, err := ioutil.ReadFile("./infrastructure/auth/jwks.json")
	if err != nil {
		return cert, err
	}

	var jwks = Jwks{}
	// err = json.NewDecoder(resp.Body).Decode(&jwks)
	err = json.Unmarshal(data, &jwks)
	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
