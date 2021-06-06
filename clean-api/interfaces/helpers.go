package interfaces

import (
	"fmt"
	"net/http"
	"strconv"
)

func getLimit(r *http.Request) uint {
	var limit uint = 5
	limitUrlQueryParam := r.URL.Query().Get("limit")
	limitU64, err := strconv.ParseUint(limitUrlQueryParam, 10, 32)
	if err != nil {
		return limit
	}
	limit = uint(limitU64)

	return limit
}

func getSkip(r *http.Request) uint {
	var skip uint = 0
	skipUrlQueryParam := r.URL.Query().Get("skip")
	skipU64, err := strconv.ParseUint(skipUrlQueryParam, 10, 32)
	if err != nil {
		return skip
	}
	skip = uint(skipU64)

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
