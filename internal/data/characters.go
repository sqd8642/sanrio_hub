package data 

import (
	"time"
	"database/sql"
	"errors"
	"github.com/lib/pq" 
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

type CharacterModel struct {
	DB *sql.DB
}

func (c CharacterModel) Insert(char *Character) error {
	query := "INSERT INTO characters (name, debut, description, hobbies, affiliations) VALUES ($1, $2, $3, $4, $5) RETURNING id, version"

	args := []any{char.Name, char.Debut, char.Description, char.Hobbies, pq.Array(char.Affiliations)}
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
	query := "UPDATE characters SET name = $1, debut = $2, description = $3, personality= $4, hobbies =$5, affiliations =$6, version = version + 1 WHERE id = $7 RETURNING version"
    
	args := []any{
		char.Name,
		char.Debut,
		char.Description,
		char.Personality,
		char.Hobbies,
		pq.Array(&char.Affiliations),
		char.ID,
	}
	
	return c.DB.QueryRow(query, args...).Scan(&char.Version)
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
	