package data 

import (
	"time"
	"database/sql"
	"errors"
	"github.com/lib/pq" 
	"context"
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
	
func (c CharacterModel) GetAll(name string, affiliations []string, filter Filters)([]*Character, error) {
	query := "SELECT id, name, debut, description, personality, hobbies, affiliations, version FROM characters WHERE (LOWER(name) = LOWER($1) OR $1 = '') AND (affiliations @> $2 OR $2 = '{}') ORDER BY id"
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel() 

	rows, err := c.DB.QueryContext(ctx, query, name, pq.Array(affiliations))
    if err != nil {
        return nil, err
    }
	defer rows.Close()

	chars := []*Character{}

	for rows.Next() {
		var char Character

		err := rows.Scan(
			&char.ID,
			&char.Name,
			&char.Debut,
			&char.Description,
			&char.Personality,
			&char.Hobbies,
			pq.Array(&char.Affiliations),
            &char.Version,
		)
		if err!= nil {
			return nil, err
		}
		chars = append(chars, &char)
	}
	if err = rows.Err(); err !=nil {
		return nil, err
	}

	return chars, nil
}