package interfaces

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type healthcheckResource struct {
	getHealthcheck application.HealthcheckService
}

type healthcheckView struct {
	Message string `json:"message"`
} //@name HealthcheckResponse

// @router /healthcheck [get]
// @summary get api and it's dependencies health status
// @tags healthcheck
// @produce  json
// @success 200 {object} common_http.InfoResponse "api response"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (hr healthcheckResource) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	err := hr.getHealthcheck.Heathcheck()
	if err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	message := "api has been eating all the apples, no need for doctors here"
	render.Render(w, r, common_http.Info(message))
}

// @router /ready [get]
// @summary get api health status
// @tags healthcheck
// @produce  json
// @success 200 {object} common_http.InfoResponse "api response"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (hr healthcheckResource) GetReady(w http.ResponseWriter, r *http.Request) {
	message := "api is up and running"
	render.Render(w, r, common_http.Info(message))
}
