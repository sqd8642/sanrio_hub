package data 

import (
	"time"
	"database/sql"
	"errors"
	"github.com/lib/pq" 
	"context"
	"fmt"
	"sanriohub.pavelkan.net/internal/validator"
	"strings"
)

type Character struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Debut time.Time `json:"debut"`
	Description string `json:"description"`
	Personality string `json:"personality"`
	Hobbies string `json:"hobbies"`
	Affiliations []string `json:"affiliations"`
	Version int32 `json:"version"`
}

func ValidateChar(v *validator.Validator, character *Character) {
	v.Check(character.Name != "","name", "must be provided")
    v.Check(character.Description != "","desc", "must be provided")
	v.Check(character.Hobbies != "","hobbies", "must be provided")
	v.Check(character.Affiliations != nil,"affiliations", "must be provided")
	v.Check(len(character.Affiliations)>=1,"affiliations", "must contain at least 1 affiliation")
	v.Check(validator.Unique(character.Affiliations), "affiliations", "must not contain duplicates")

}

type CharacterModel struct {
	DB *sql.DB
}

func (c CharacterModel) Insert(char *Character) error {
	query := "INSERT INTO characters (name, debut, description, hobbies, affiliations, personality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, version"

	args := []any{char.Name, char.Debut, char.Description, char.Hobbies, pq.Array(char.Affiliations), char.Personality}
    return c.DB.QueryRow(query, args...).Scan(&char.ID, &char.Version)

}

func (c CharacterModel) Get(id int64) (*Character, error){
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := "SELECT id, name, debut, description, personality, hobbies, affiliations, version FROM characters WHERE id = $1"
	var char Character 

	err:= c.DB.QueryRow(query, id).Scan(
		&char.ID,
		&char.Name,
	    &char.Debut,
		&char.Description,
		&char.Personality,
		&char.Hobbies,
		pq.Array(&char.Affiliations),
		&char.Version,
	)

	if err != nil {
		switch {
		    case errors.Is(err, sql.ErrNoRows):
		        return nil, ErrRecordNotFound
		    default:
		        return nil, err
		    }
		}
		return &char, nil		
}

func (c CharacterModel) Update(char *Character) error {
	query := "UPDATE characters SET name = $1, debut = $2, description = $3, personality= $4, hobbies =$5, affiliations =$6, version = version + 1 WHERE id = $7 AND version = $8 RETURNING version"
    
	args := []any{
		char.Name,
		char.Debut,
		char.Description,
		char.Personality,
		char.Hobbies,
		pq.Array(&char.Affiliations),
		char.ID,
		char.Version,
	}
	
	err := c.DB.QueryRow(query, args...).Scan(&char.Version)
	if err != nil {
		switch{
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (c CharacterModel) Delete(id int64) error {
	if id <1 {
		return ErrRecordNotFound
	}

	query := "DELETE FROM characters where id = $1"

	result, err := c.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
	
func (c CharacterModel) GetAll(name string, affiliations []string, filter Filters)([]*Character, Metadata, error) {
	query :=  fmt.Sprintf("SELECT count(*) OVER(), id, name, debut, description, personality, hobbies, affiliations, version FROM characters WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') AND (affiliations @> $2 OR $2 = '{}') ORDER BY %s %s, id ASC LIMIT $3 OFFSET $4", filter.sortColumn(), filter.sortDirection())
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel() 

	rows, err := c.DB.QueryContext(ctx, query, name, pq.Array(affiliations),  filter.limit(), filter.offset() )
    if err != nil {
		return nil, Metadata{}, err 
	}
		
	defer rows.Close()

	totalRecords := 0
	chars := []*Character{}

	for rows.Next() {
		var char Character

		err := rows.Scan(
			&totalRecords,
			&char.ID,
			&char.Name,
			&char.Debut,
			&char.Description,
			&char.Personality,
			&char.Hobbies,
			pq.Array(&char.Affiliations),
            &char.Version,
		)
		if err != nil {
			return nil, Metadata{}, err 
		}
			
		chars = append(chars, &char)
	}
	if err = rows.Err(); err !=nil {
		return nil,Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return chars, metadata, nil
}

func (c CharacterModel) GetCharactersByID(ids []int64) ([]*Character, error) {
    // Prepare the placeholders for the IN clause
    placeholders := make([]string, len(ids))
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        placeholders[i] = fmt.Sprintf("$%d", i+1)
        args[i] = id
    }
    // Construct the IN clause
    inClause := strings.Join(placeholders, ", ")

    // Construct the SQL query with the IN clause
    query := fmt.Sprintf(`
        SELECT id, name
        FROM characters
        WHERE id IN (%s);
    `, inClause)

    // Execute the query
    rows, err := c.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Iterate over the rows and scan characters
    var characters []*Character
    for rows.Next() {
        var char Character
        if err := rows.Scan(&char.ID, &char.Name); err != nil {
            return nil, err
        }
        characters = append(characters, &char)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return characters, nil
}