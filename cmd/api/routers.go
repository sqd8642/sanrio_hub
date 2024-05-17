package main 

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

    router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/characters", app.addCharHandler) 
    router.HandlerFunc(http.MethodGet, "/v1/characters/:id",app.showCharHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/characters/:id", app.updateCharHandler)
	router.HandlerFunc(http.MethodDelete,"/v1/characters/:id", app.requirePermission("characters:write", app.deleteCharHandler))
	router.HandlerFunc(http.MethodGet,"/v1/characters", app.listCharsHandler)

	router.HandlerFunc(http.MethodPost, "/v1/shows", app.addShowHandler)
    router.HandlerFunc(http.MethodGet, "/v1/shows/:id", app.showShowHandler)
    router.HandlerFunc(http.MethodPatch, "/v1/shows/:id", app.updateShowHandler)
    router.HandlerFunc(http.MethodDelete, "/v1/shows/:id", app.requirePermission("show:write", app.deleteShowHandler))
    router.HandlerFunc(http.MethodGet, "/v1/shows", app.listShowsHandler)

	router.HandlerFunc(http.MethodGet, "/v1/show/:id/characters", app.listShowCharsHandler)


	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
    
	return app.authenticate(app.authenticate(router))
}
