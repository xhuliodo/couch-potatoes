package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/form3tech-oss/jwt-go"
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/logger"
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

	data, err := ioutil.ReadFile("/auth/jwks.json")
	if err != nil {
		os.Exit(0)
	}

	err = json.Unmarshal(data, &jwks)
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

func CacheJwksCert(errorLogger *logger.ErrorLogger) {
	jwks := Jwks{}

	resp, err := http.Get(os.Getenv("JWKS_URL"))
	if err != nil {
		cause := errors.New("check jwks url, the one provided does send back any response")
		errStack := errors.Wrap(cause, err.Error())
		errorLogger.Log(errStack)
		log.Fatal(errStack)
		os.Exit(0)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		cause := errors.New("check jwks json format, the one provided is not standard")
		errStack := errors.Wrap(cause, err.Error())
		errorLogger.Log(errStack)
		log.Fatal(errStack)
		os.Exit(0)
	}

	certFromOnline, err := json.MarshalIndent(jwks, "", " ")
	if err != nil {
		cause := errors.New("could not marshall the json")
		errStack := errors.Wrap(cause, err.Error())
		errorLogger.Log(errStack)
		log.Fatal(errStack)
		os.Exit(0)
	}

	err = ioutil.WriteFile("/auth/jwks.json", certFromOnline, os.ModePerm)
	if err != nil {
		cause := errors.New("could not cache cert to container")
		errStack := errors.Wrap(cause, err.Error())
		errorLogger.Log(errStack)
		log.Fatal(errStack)
		os.Exit(0)
	}
}
