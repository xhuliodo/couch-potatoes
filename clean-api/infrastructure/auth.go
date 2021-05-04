package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
		pemCert := []byte("-----BEGIN CERTIFICATE-----\nMIIDDTCCAfWgAwIBAgIJK1eSJpVnlXzQMA0GCSqGSIb3DQEBCwUAMCQxIjAgBgNV\nBAMTGWRldi1wczVkcXFpcy5ldS5hdXRoMC5jb20wHhcNMjEwMTIyMTExODQwWhcN\nMzQxMDAxMTExODQwWjAkMSIwIAYDVQQDExlkZXYtcHM1ZHFxaXMuZXUuYXV0aDAu\nY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1XPG2/iq3ET6UkwL\nahk+VXpbPfL098831+EcMfsFBtWTSwQPRQBXfeE6o/dVmt3nI46Ddi2xk7siQV1r\nAFD2IliLsoPJHTMGSc1VHkB8qhcLcyN1kW9z/J6LPhdW1WOEhANZvXKsyfbX3lN6\nd04LgRGjnJz5K1G5S8Ba9ggDvHQnWF/THTjTFTPRiPr5H6EgJTzVA2uSqO1QFIKF\nBUTALfCOPLQeb6GyovE2EH0r7ASZsNVss4htKAi3h57Kh1GFH5nOPlwVfK7O0+af\ndVKQdOV52Oyr7IMdXaX3uOV0/aSst0MoWVJFb3m1HhYwCH6AUliWviJg3/TqdzJt\nRzvMOwIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQ3BX7/pham\nENNFaAtcFvmb6F4rDzAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEB\nAG1i8y+/fSbyBBk5GSKRSaEj0NKdPqAfgaIFcPg15IYs1i1Pjy152ReZ3C4ksrGn\nn2TRORtl8JAB+xKMkBAgFl7gdL7SxhoPGMNFq4xf2DkSrQ5vCHYa8OadzQ+2ij8r\ncn3ngP3KY9Byjo6kGizmHSPo0oO74UBYq6D4tCpeQkVRKu6n9ZJ0vy2t/eQXa5ww\ngSFqdgVbMKmI+lFj0FKLSzdf/vpQEfhMBUnm2aHN2xef7fp4ezSoHcuJ3PYscgKB\nD1CitfxdD0g/Qs47FAdLeQvjRveXc4vn2qqv5Np9RbQo9PNH1Iy/gH6o6bUOmex4\nTy59tzK2cIcXSJLDSICgTKs=\n-----END CERTIFICATE-----")
		cert, _ := jwt.ParseRSAPublicKeyFromPEM(pemCert)
		return cert, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
	Extractor:     jwtmiddleware.FromAuthHeader,
	// TODO: implement custom unauthorized error message
	// ErrorHandler:  onAuthorizedReq,
})

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := jwtMiddleware.CheckJWT(rw, r)
		if err != nil {
			// _ = render.Render(rw, r, common_http.ErrBadRequest(err))
			return
		}

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

// func onAuthorizedReq(w http.ResponseWriter, r *http.Request, err error) {
// 	_ = render.Render(w, r, common_http.ErrBadRequest(err))

// }
