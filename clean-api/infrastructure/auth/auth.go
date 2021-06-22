package auth

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

func onUnauthorizedReq(w http.ResponseWriter, r *http.Request, errString string) {
	cause := errors.New("not_authenticated")
	errStack := errors.Wrap(cause, "please provide a valid jwt token")
	render.Render(w, r, common_http.DetermineErr(errStack))
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("authorization")
		encodedToken = strings.Replace(encodedToken, "Bearer ", "", -1)
		decodedToken, _ := jwt.Parse(encodedToken, nil)

		if err := validateToken(decodedToken, encodedToken); err != nil {
			onUnauthorizedReq(rw, r, err.Error())
			return
		}

		c := r.Context()

		tokenClaims := decodedToken.Claims.(jwt.MapClaims)
		userId := tokenClaims["sub"]
		isAdmin := getIsAdmin(tokenClaims)

		rWithUserId := r.WithContext(context.WithValue(c, "userId", userId))
		rWithAdmin := r.WithContext(context.WithValue(rWithUserId.Context(), "isAdmin", isAdmin))

		*r = *rWithAdmin

		next.ServeHTTP(rw, r)
	})
}

func validateToken(token *jwt.Token, stringToken string) error {
	_, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		pemCert, err := GetPemCert(token)
		if err != nil {
			return nil, err
		}
		cert, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pemCert))
		if err != nil {
			return nil, err
		}
		return cert, nil
	})
	
	if err != nil {
		if err != nil {
			e, ok := err.(*jwt.ValidationError)
			if !ok || ok && e.Errors&jwt.ValidationErrorIssuedAt == 0 { // Don't report error that token used before issued.
				return err
			}
		}
	}

	return nil
}

func getIsAdmin(claims jwt.MapClaims) bool {
	isAdmin := claims["https://couch-potatoes.com/claims/"]
	if isAdmin != nil {
		isAdminArray := reflect.ValueOf(isAdmin)
		isAdminInterface := isAdminArray.Index(0)

		isAdminString := fmt.Sprintf("%v", isAdminInterface)

		return isAdminString == "admin"
	}
	return false

}
