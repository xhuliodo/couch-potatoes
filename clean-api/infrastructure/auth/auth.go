package auth

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
		pemCert, err := GetPemCert(t)
		if err != nil {
			panic(err.Error())
		}
		cert, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(pemCert))
		return cert, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
	ErrorHandler:  onUnauthorizedReq,
})

func onUnauthorizedReq(w http.ResponseWriter, r *http.Request, errString string) {
	cause := errors.New("not_authenticated")
	errStack := errors.Wrap(cause, "please provide a valid jwt token")
	render.Render(w, r, common_http.DetermineErr(errStack))
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jwtMiddleware.CheckJWT(rw, r)

		c := r.Context()
		stringToken := c.Value(jwtMiddleware.Options.UserProperty)
		token := stringToken.(*jwt.Token)
		tokenClaims := token.Claims.(jwt.MapClaims)
		userId := tokenClaims["sub"]
		isAdmin := getIsAdmin(tokenClaims)

		rWithUserId := r.WithContext(context.WithValue(c, "userId", userId))
		rWithAdmin := r.WithContext(context.WithValue(rWithUserId.Context(), "isAdmin", isAdmin))

		*r = *rWithAdmin

		next.ServeHTTP(rw, r)
	})
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
