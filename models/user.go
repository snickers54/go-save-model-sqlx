package models
import (
    _ "log"
    "database/sql"
)
type User struct {
    AutoIncr
    Username sql.NullString `db:"username" json:"username"`
    FacebookID string `db:"facebook_id" json:"facebook_id"`
    GoogleID string `db:"google_id" json:"google_id"`
    Email string `db:"email" json:"email"`
    Password string `db:"password" json:"-"`
    Test bool `db:"test", json:"test"`
}

type Users []User

func (u *User) Table() string {
    return "user"
}

// example of implementation of Save method for my custom model

func (u *User) Save() error {
    var err error
    var lastId int64
    if u.ID > 0 {
        _, err = Update(u, GetInstance())
    } else {
        lastId, err = Save(u, GetInstance())
        u.ID = uint64(lastId)
    }
    return err
}

func (u *User) GetById(id string) error {
    return u.db().Get(u, "SELECT * FROM user WHERE id = '?'", id)
}

func (u *User) DeleteById(id string) error{
    _, err := u.db().Exec("DELETE FROM user WHERE id = ?", id)
    return err
}

func (u *User) Get(limit, offset, order, sense string) (Users, error) {
    users := Users{}
    err := u.db().Get(&users, "SELECT * FROM user WHERE limit ?,? ORDER BY ? ?", limit, offset, order, sense)
    return users, err
}
