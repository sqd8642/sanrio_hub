package main 

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

    router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/characters",  app.requirePermission("characters:read", app.addCharHandler)) 
    router.HandlerFunc(http.MethodGet, "/v1/characters/:id",  app.requirePermission("characters:read",app.showCharHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/characters/:id",  app.requirePermission("characters:read", app.updateCharHandler))
	router.HandlerFunc(http.MethodDelete,"/v1/characters/:id", app.requirePermission("characters:read", app.deleteCharHandler))
	router.HandlerFunc(http.MethodGet,"/v1/characters", app.requirePermission("characters:read",app.listCharsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
    
	return app.authenticate(app.authenticate(router))
}
