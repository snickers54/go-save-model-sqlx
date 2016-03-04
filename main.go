package main

import (
    "github.com/snickers54/go-save-model-sqlx/models"
    "database/sql"
)

func main() {
    u := models.User{}
    u.Username = sql.NullString{}
    u.Username.String = "Toto"
    u.Username.Valid = true
    u.Save()
}
