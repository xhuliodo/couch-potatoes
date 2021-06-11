package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	var jwks = Jwks{}

	data, err := ioutil.ReadFile("./infrastructure/auth/jwks.json")
	if err != nil {
		err = getCertPemFromOnline(&jwks)
		if err != nil {
			return cert, err
		}

		err = cacheCertPem(jwks)
		if err != nil {
			return cert, err
		}
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, &jwks)
		if err!=nil{
			return cert, err
		}
	}

	if err != nil {
		fmt.Println("if fucks up even though you get the cert online, with err: ", err)
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

func getCertPemFromOnline(jwks *Jwks) (err error) {
	resp, err := http.Get("https://dev-ps5dqqis.eu.auth0.com/.well-known/jwks.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return err
	}

	return nil
}

func cacheCertPem(jwks Jwks) error {
	certFromOnline, _ := json.Marshal(jwks)
	err := ioutil.WriteFile("./infrastructure/auth/jwks.json", certFromOnline, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
