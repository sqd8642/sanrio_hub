package main
import (
    "fmt"
	"time"
    "net/http"
    "sanriohub.pavelkan.net/internal/data"
    
)

func (app *application) addCharHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int64 `json:"id"`
	    Name string `json:"name"`
	    Debut time.Time `json:"debut"`
	    Description string `json:"description"`
	    Personality string `json:"personality"`
	    Hobbies string `json:"hobbies"`
	    Affiliations []string `json:"affiliations"`
	    Version int32 `json:"version"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err.Error())
        return
	}
    fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showCharHandler(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, any interpolated URL parameters will be
	// stored in the request context. We can use the ParamsFromContext() function to
	// retrieve a slice containing these parameter names and values.
	
	id, err := app.readIDParam(r)
	
	if err != nil {
	    http.NotFound(w, r)
	    return
	} 

	character := data.Character{
		ID: id,
	    Name: "Kuromi", 
	    Debut: time.Date(2005, time.September, 1, 0, 0, 0, 0, time.UTC),
	    Description: "Kuromi is a character from the My Melody universe. She is My Melody's rival and doppelgänger, and manifests as a white rabbit or imp-like creature wearing a black jester's hat with a pink skull on the front and a black devil's tail. The skull's facial expression on her forehead changes to match Kuromi's mood. Fittingly, her birthday is on Halloween (October 31st).",
	    Personality: "Although Kuromi may look and act tough and punk, she is actually very girly and is attracted to good-looking guys! Kuromi enjoys writing in her diary and is hooked on romantic short stories. Her favorite colors are black and hot pink. Her favorite food is shallots, all kinds of meat, and in recent Kuromi merchandise, cherries have been shown. Despite being a villain (in Onegai My Melody), Kuromi is mostly into food and even cooks. You might call her a rowdy free spirit.Kuromi is the punk tomboy counterpart to My Melody. Although My Melody can get along with Kuromi, the latter's feelings are more inclined to rivalry in order to appear tough.",
	    Hobbies: "Kuromi's hobbies include writing in his diary, cooking and reading romance novels.",
	    Affiliations: []string{"Hello Kitty", "MyMelody", "Pompompurin", "Chococat", "Keroppi", "Cinnamoroll", "Pochacco", "Badtz-Maru"},
	    Version: 1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
    if err != nil { 
        app.logger.Print(err)
        http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
    }
}

	
	
