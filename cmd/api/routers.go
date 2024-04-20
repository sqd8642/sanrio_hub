package main 

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()

    router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/characters", app.addCharHandler)
    router.HandlerFunc(http.MethodGet, "/v1/characters/:id", app.showCharHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/characters/:id", app.updateCharHandler)
	router.HandlerFunc(http.MethodDelete,"/v1/characters/:id", app.deleteCharHandler)
	router.HandlerFunc(http.MethodGet,"/v1/characters", app.listCharsHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)


	return router
}
