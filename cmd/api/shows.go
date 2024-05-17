package main

import (
	"errors"
	"fmt"
	"net/http"
	"sanriohub.pavelkan.net/internal/data"
	"sanriohub.pavelkan.net/internal/validator"
	"time"
)

func (app *application) addShowHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string    `json:"title"`
		ReleaseDate time.Time `json:"release_date"`
		Description string    `json:"description"`
		Genre       string    `json:"genre"`
		Duration    int       `json:"duration"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	show := &data.Show{
		Title:       input.Title,
		ReleaseDate: input.ReleaseDate,
		Description: input.Description,
		Genre:       input.Genre,
		Duration:    input.Duration,
		// Initialize other fields
	}

	v := validator.New()

	if data.ValidateShow(v, show); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Shows.Insert(show)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/shows/%d", show.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"show": show}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showShowHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	show, err := app.models.Shows.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"show": show}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateShowHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    show, err := app.models.Shows.Get(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            app.notFoundResponse(w, r)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    var input struct {
        Title       *string   `json:"title"`
        ReleaseDate *time.Time `json:"release_date"`
        Description *string   `json:"description"`
        Genre       *string   `json:"genre"`
        Duration    *int      `json:"duration"`
        // Add other fields as needed
    }

    err = app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    if input.Title != nil {
        show.Title = *input.Title
    }
    if input.ReleaseDate != nil {
        show.ReleaseDate = *input.ReleaseDate
    }
    if input.Description != nil {
        show.Description = *input.Description
    }
    if input.Genre != nil {
        show.Genre = *input.Genre
    }
    if input.Duration != nil {
        show.Duration = *input.Duration
    }
    // Update other fields similarly

    v := validator.New()

    if data.ValidateShow(v, show); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = app.models.Shows.Update(show)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"show": show}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}

func (app *application) deleteShowHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    err = app.models.Shows.Delete(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            app.notFoundResponse(w, r)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"message": "show successfully deleted"}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}

func (app *application) listShowsHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title     string
        Genre     string
        Filters   data.Filters
    }

    v := validator.New()
    qs := r.URL.Query()

    input.Title = app.readString(qs, "title", "")
    input.Genre = app.readString(qs, "genre", "")
    input.Filters.Page = app.readInt(qs, "page", 1, v)
    input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
    input.Filters.Sort = app.readString(qs, "sort", "id")
    input.Filters.SortSafelist = []string{"id", "title", "release_date", "-id", "-title", "-release_date"}

    if data.ValidateFilters(v, input.Filters); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    shows, metadata, err := app.models.Shows.GetAll(input.Title, input.Genre, input.Filters)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"shows": shows, "metadata": metadata}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}


func (app *application) listShowCharsHandler( w http.ResponseWriter, r *http.Request) {
    showID, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    ids, err := app.models.ShowCharacters.GetCharactersByShowID(showID)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    characters, err := app.models.Characters.GetCharactersByID(ids)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
    err = app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}