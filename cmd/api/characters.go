package main
import (
    "fmt"
	"time"
    "net/http"
    "sanriohub.pavelkan.net/internal/data"
	"sanriohub.pavelkan.net/internal/validator"
	"errors"
    
)

func (app *application) addCharHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
	    Name string `json:"name"`
	    Debut time.Time `json:"debut"`
	    Description string `json:"description"`
	    Personality string `json:"personality"`
	    Hobbies string `json:"hobbies"`
	    Affiliations []string `json:"affiliations"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
		
    char := &data.Character{
		Name: input.Name,
	    Debut: input.Debut,
	    Description: input.Description,
	    Personality: input.Personality,
	    Hobbies: input.Hobbies,
	    Affiliations: input.Affiliations,
	}

	v := validator.New()

	
	if data.ValidateChar(v, char); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

    err = app.models.Characters.Insert(char)
    if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("/v1/movies/%d", char.ID))

    err = app.writeJSON(w, http.StatusCreated, envelope{"char": char}, headers)
    if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showCharHandler(w http.ResponseWriter, r *http.Request) {
	
	id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }
    
	char, err := app.models.Characters.Get(id)
	if err != nil {
		switch {
		    case errors.Is(err, data.ErrRecordNotFound):
		        app.notFoundResponse(w, r)
		    default:
		        app.serverErrorResponse(w, r, err)
		}
		return
	}
		
	err = app.writeJSON(w, http.StatusOK, envelope{"character": char}, nil)
    if err != nil {
		app.serverErrorResponse(w, r, err)
	}	
}

func (app *application) updateCharHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

	char, err := app.models.Characters.Get(id)
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
		Name *string `json: "title"`
		Debut *time.Time  `json: "debut"`
		Description *string `json:"description"`
	    Personality *string `json:"personality"`
	    Hobbies *string `json:"hobbies"`
	    Affiliations []string `json:"affiliations"`
	}

	err = app.readJSON(w, r, &input)

    if err != nil {
        app.badRequestResponse(w, r, err)
        return
	}

    if input.Name != nil {
		char.Name = *input.Name
	}
	if input.Debut != nil {
		char.Debut = *input.Debut
	}
	if input.Description != nil {
		char.Description = *input.Description
	}
	if input.Personality != nil {
		char.Personality = *input.Personality
	}
	if input.Hobbies != nil {
		char.Hobbies = *input.Hobbies
	}
	if input.Affiliations != nil {
		char.Affiliations = input.Affiliations
	}
	v := validator.New()

	if data.ValidateChar(v, char); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Characters.Update(char)
	if err != nil {
        switch {
            case errors.Is(err, data.ErrEditConflict):
                app.EditConflictResponse(w, r)
            default:
                app.serverErrorResponse(w, r, err)
        }
        return
    }


	err = app.writeJSON(w, http.StatusOK, envelope{"character":char}, nil)
	if err!= nil{
		app.serverErrorResponse(w,r,err)
	}		
}

func (app *application) deleteCharHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w,r)
		return
	}

	err = app.models.Characters.Delete(id)
	if err!= nil {
		switch  {
	        case errors.Is(err, data.ErrRecordNotFound):
		        app.notFoundResponse(w, r)
		    default:
		        app.serverErrorResponse(w, r, err)
	    }
	    return
    }
		err = app.writeJSON(w, http.StatusOK, envelope{"message": "character successfully deleted"}, nil)
		if err != nil {
		app.serverErrorResponse(w, r, err)
		}
}

func (app *application) listCharsHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string 
		Affiliations []string 
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
    input.Affiliations = app.readCSV(qs, "affiliations", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
    input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
    input.Filters.SortSafelist = []string{"id", "name", "debyt", "-id", "-name", "-debut"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	    return
	}
	chars, metadata, err := app.models.Characters.GetAll(input.Name, input.Affiliations, input.Filters)
	if err!= nil {
		app.serverErrorResponse(w,r,err)
		return 
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"characters": chars, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w,r, err)
	}
}

