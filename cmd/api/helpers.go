package main
import (
    "encoding/json" 
    "net/http"
	"errors"
    "strconv"
    "net/url"
	"github.com/julienschmidt/httprouter"
    "fmt" 
    "io" 
    "strings"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

    js, err := json.MarshalIndent(data, "", "\t")

    if err != nil {
        return err
    }

    js = append(js, '\n')

    for key, value := range headers {
        w.Header()[key] = value
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(js)

    return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
	    return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
  
    err := json.NewDecoder(r.Body).Decode(dst)
    if err != nil {
 
    var syntaxError *json.SyntaxError
    var unmarshalTypeError *json.UnmarshalTypeError
    var invalidUnmarshalError *json.InvalidUnmarshalError
    switch {
    
        case errors.As(err, &syntaxError):
            return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
   
        case errors.Is(err, io.ErrUnexpectedEOF):
            return errors.New("body contains badly-formed JSON")
    
        case errors.As(err, &unmarshalTypeError):
            if unmarshalTypeError.Field != "" {
                return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
            }
            return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
    
        case errors.Is(err, io.EOF):
            return errors.New("body must not be empty")
        case errors.As(err, &invalidUnmarshalError):
    panic(err)
   
        default:
            return err
            }
        }
        return nil
    }

    func (app * application) readString (qs url.Values, key string, defaultValue string) string {
        s:=qs.Get(key)

        if s == "" {
            return defaultValue
        }
        return s
    }

    func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
        csv := qs.Get(key)

        if csv == "" {
            return defaultValue
        }
        return strings.Split(csv, ",")
    }

    func (app *application) readInt(qs url.Values, key string, defaultValue int) int {
        s:=qs.Get(key)

        if s == ""{
            return defaultValue
        }
        i, err := strconv.Atoi(s)
    
        if err != nil {
            return defaultValue
        }
        
        return i
    }