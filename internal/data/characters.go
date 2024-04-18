package data 

import (
	"time"
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