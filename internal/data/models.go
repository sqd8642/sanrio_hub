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
    Shows ShowModel
    Tokens TokenModel
    Permissions PermissionModel
    Users UserModel
    ShowCharacters ShowCharactersModel
}

func NewModels(db *sql.DB) Models {
    return Models{
        Characters: CharacterModel{DB: db},
        Shows: ShowModel{DB: db},
        Tokens: TokenModel{DB:db},
        Permissions: PermissionModel{DB: db},
        Users: UserModel{DB: db},
        ShowCharacters: ShowCharactersModel{DB: db},
    }
}