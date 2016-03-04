package models

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "log"
)

var _instance *sqlx.DB = nil

func GetInstance() *sqlx.DB {
    if _instance == nil {
        connect()
    }
    return _instance
}

func connect() bool {
    db, err := sqlx.Connect("mysql", "root:toor@/test")
    if err != nil {
        log.Println(err)
        return false
    }
    _instance = db
    return true
}
