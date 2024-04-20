package data

import (
    "database/sql"
    "errors"
)

var (
    ErrRecordNotFound = errors.New("record not found")
    ErrEditConflict = errors.New("edit conflict")

)

type Models struct {
    Characters CharacterModel
    Tokens TokenModel
    Users UserModel
}

func NewModels(db *sql.DB) Models {
    return Models{
        Characters: CharacterModel{DB: db},
        Tokens: TokenModel{DB:db},
        Users: UserModel{DB: db},
    }
}