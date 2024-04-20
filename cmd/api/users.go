package main
import (
    "time"
	"errors"
    "net/http"
    "sanriohub.pavelkan.net/internal/data"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

    var input struct {
        Name string `json:"name"`
        Email string `json:"email"`
        Password string `json:"password"`
    }
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
    return
    }
    user := &data.User{
    Name: input.Name,
    Email: input.Email,
    Activated: false,
    }
    err = user.Password.Set(input.Password)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    } 
    err = app.models.Users.Insert(user)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)

	data := map[string]any{
		"activationToken": token.Plaintext,
		"userID": user.ID,
		}
		
    err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)

	if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        TokenPlaintext string `json:"token"`
    }
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }


    user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
    if err != nil {
            app.serverErrorResponse(w, r, err)
        return
    }

    user.Activated = true
    err = app.models.Users.Update(user)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrEditConflict):
            app.EditConflictResponse(w, r)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
