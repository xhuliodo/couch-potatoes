package interfaces

import (
	"fmt"
	"net/http"
	"strconv"
)

// A ResourceCreatedView is an response that is used when the resource creation was successful
// swagger:response resourceCreatedView
type ResourceCreatedView struct {
	Message string `json:"message"`
}

const (
	defaultLimit uint = 5
	defaultSkip  uint = 0
)

func getLimit(r *http.Request) uint {
	limitUrlQueryParam := r.URL.Query().Get("limit")
	limitU64, err := strconv.ParseUint(limitUrlQueryParam, 10, 32)
	if err != nil {
		return defaultLimit
	}
	limit := uint(limitU64)

	return limit
}

func getSkip(r *http.Request) uint {
	skipUrlQueryParam := r.URL.Query().Get("skip")
	skipU64, err := strconv.ParseUint(skipUrlQueryParam, 10, 32)
	if err != nil {
		return defaultSkip
	}
	skip := uint(skipU64)

	return skip
}

func getUserId(r *http.Request) string {
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	return userId
}

func getIsAdmin(r *http.Request) bool {
	isAdminInterface := r.Context().Value("isAdmin")
	isAdmin := isAdminInterface.(bool)

	return isAdmin
}
