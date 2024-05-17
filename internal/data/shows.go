package data

import (
	"context"
	"database/sql"
	"fmt"

	"time"
	"errors"

	"sanriohub.pavelkan.net/internal/validator"
)

type Show struct {
	ID           string    `json:"id"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
	Title        string    `json:"title"`
	ReleaseDate  time.Time `json:"releaseDate"`
	Description  string    `json:"description"`
	Genre        string    `json:"genre"`
	Duration     int       `json:"duration"`
}

type ShowModel struct {
	DB       *sql.DB
}

func ValidateShow(v *validator.Validator, show *Show) {
    v.Check(show.Title != "", "title", "must be provided")
    v.Check(show.Description != "", "description", "must be provided")
    v.Check(show.Genre != "", "genre", "must be provided")
    v.Check(show.Duration > 0, "duration", "must be greater than 0")
}

func (m ShowModel) GetAll(title string, genre string,  filters Filters) ([]*Show, Metadata, error) {
	// Retrieve all shows from the database.
	query := fmt.Sprintf("SELECT count(*) OVER(), id, created_at, updated_at,title, release_date, description, genre, duration FROM shows WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') AND (to_tsvector('simple', genre) @@ plainto_tsquery('simple', $2) OR $2 = '') ORDER BY %s %s, id ASC LIMIT $3 OFFSET $4", filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{title, genre, filters.limit(), filters.offset()}

	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer rows.Close()

	// Declare a totalRecords variable
	totalRecords := 0

	var shows []*Show
	for rows.Next() {
		var show Show
		err := rows.Scan(&totalRecords, &show.ID, &show.CreatedAt, &show.UpdatedAt, &show.Title, &show.ReleaseDate, &show.Description, &show.Genre, &show.Duration)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Show struct to the slice
		shows = append(shows, &show)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the shows and metadata.
	return shows, metadata, nil
}

func (m ShowModel) Insert(show *Show) error {
	query := `
		INSERT INTO shows (title, release_date, description, genre, duration) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
	`
	args := []interface{}{show.Title, show.ReleaseDate, show.Description, show.Genre, show.Duration}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&show.ID, &show.CreatedAt, &show.UpdatedAt)
}

func (m ShowModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM shows
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m ShowModel) Update(show *Show) error {
	query := `
		UPDATE shows
		SET title = $1, release_date = $2, description = $3, genre = $4, duration = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND updated_at = $7
		RETURNING updated_at
	`
	args := []interface{}{show.Title, show.ReleaseDate, show.Description, show.Genre, show.Duration, show.ID, show.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&show.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m ShowModel) Get(id int64) (*Show, error) {
    // Retrieve a specific show based on its ID.
    query := `
        SELECT id, created_at, updated_at, title, release_date, description, genre, duration
        FROM shows
        WHERE id = $1
        `
    var show Show
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    row := m.DB.QueryRowContext(ctx, query, id)
    err := row.Scan(&show.ID, &show.CreatedAt, &show.UpdatedAt, &show.Title, &show.ReleaseDate, &show.Description, &show.Genre, &show.Duration)
    if err != nil {
        return nil, fmt.Errorf("cannot retrieve show with id %s: %w", id, err)
    }
    return &show, nil
}

type ShowCharacter struct {
    ID          int    `json:"id"`
    ShowID      int    `json:"show_id"`
    CharacterID int    `json:"character_id"`
    // Add other fields as needed
}

// ShowCharactersModel holds methods for querying the ShowCharacters table.
type ShowCharactersModel struct {
    DB *sql.DB
}


func (m ShowCharactersModel) GetCharactersByShowID(showID int64) ([]int64, error) {
    query := `
        SELECT character_id
        FROM ShowCharacters
        WHERE show_id = $1
    `
    rows, err := m.DB.Query(query, showID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ids []int64
    for rows.Next() {
        var id int64
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        ids = append(ids, id)
    }
    
    return ids, nil
}
