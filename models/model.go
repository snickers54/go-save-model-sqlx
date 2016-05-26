package models

import (
    "time"
    "github.com/jmoiron/sqlx"
    "database/sql/driver"
    "fmt"
    "reflect"
    "strings"
    "log"
)
type AutoIncr struct {
    ID       uint64 `db:"id" json:"id" primary:"-"`
    Created  time.Time `db:"created_at" json:"created_at"`
    Updated  time.Time `db:"updated_at" json:"updated_at"`
    // this will not be used during the reflection
}

type Model interface {
    Save() error
}

type Table interface {
    Table() string
}

func bool2int(b bool) int {
   if b {
      return 1
   }
   return 0
}

func reflectStatements(obj interface{}, withPrimary bool) string {
    // we reflect via ValueOf because we want to access dynamic values
    val := reflect.ValueOf(obj)
    // if a pointer to a struct is passed, get the vale of the dereferenced object
    if val.Kind() == reflect.Ptr {
      val = val.Elem()
    }

    var str []string
    // if what is given is not a structure, we just return empty string and it wont be counted
    if val.Kind() != reflect.Struct {
        return ""
    }
    // we will need this later to check if we are not in presence of types like sql.NullString
    valuerType := reflect.TypeOf((*driver.Valuer)(nil)).Elem()
    stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
    // we iterate on our struct fields
    
    whereStr := ""

    for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
        // if we detect this field has a tag called `primary` and we don't want primary key to be used (for INSERT) we continue
        if len(tag.Get("primary")) > 0 && withPrimary == false {
	   continue;
	} else if len(tag.Get("primary")) > 0 && withPrimary == true {
	  var r interface{} = valueField.Interface()
	  whereStr = " WHERE " + fmt.Sprintf("`%s` = '%v'", tag.Get("db"), r)
	  continue;
	}
        // log.Println(tag.Get("db"), valueField.Kind().String(), valueField.MethodByName("String").IsValid(), valueField.Type())
        // Here, if we detect the field itself is a struct, and if this struct doesn't implement neither the stringer and valuer interfacep we call ourselves
        if valueField.Kind().String() == "struct" && valueField.Type().Implements(stringerType) == false &&
            valueField.Type().Implements(valuerType) == false {
            if recursive := reflectStatements(valueField.Interface(), withPrimary); len(recursive) > 0 {
                str = append(str, recursive)
            }
            continue
        }
        // if there is no tag db OR the tag db is explicitly telling us we don't want to use this field, we continue
        if tag.Get("db") == "" || tag.Get("db") == "-" {continue;}
        // this part is handling `options` types like sql.NullString, sql.NullInt64 ...
        // we get the struct and put it into a generic interface
        var rawValue interface{} = valueField.Interface()
        // if this struct implement the interface Valuer, it means it's a type like listed above
        if valueField.Type().Implements(valuerType) {
            var err error
            // so, we get the Value and rewrite our generic inteface with it
            rawValue, err = valueField.Interface().(driver.Valuer).Value()
            // log.Println(err)
            // if there is either an error we don't include it
            if err != nil {continue}
        }
        var value string
        // here we treat our value to be concatenate to our statement
        // if it's a string OR a sql.NullString (didn't find a better way for now) it's dquoted encapsulated
        if valueField.Kind().String() == "string" ||
            valueField.Type().Implements(stringerType) ||
            valueField.Type().String() == "sql.NullString" {
            value = fmt.Sprintf("\"%s\"", rawValue)
        // if it's nil, we put null
        } else if rawValue == nil {
            value = fmt.Sprintf("%s", "NULL")
        // if it's a boolean, we have to convert it to an int, because sql is handling boolean with a tiny int right ?!
        } else if valueField.Kind().String() == "bool" {
            value = fmt.Sprintf("%d", bool2int(rawValue.(bool)))
        // if we don't know what it is, my best guess is to let %v find the best value to concatenate
        } else {
            value = fmt.Sprintf("%v", rawValue)
        }
        // we now append it to key = value array of strings
        str = append(str, fmt.Sprintf("`%s` = %s", tag.Get("db"), value))
	}
    // log.Println(str)
    // we simply join them with a coma
    return strings.Join(str, ",") + whereStr
}

func getTable(obj interface{}) string {
    val := reflect.ValueOf(obj)
    // if a pointer to a struct is passed, get the vale of the dereferenced object
    if val.Kind() == reflect.Ptr {
      val = val.Elem()
    }
    if _, found := obj.(Table); found {
        return obj.(Table).Table()
    }
    return strings.ToLower(val.Type().Name())
}

func Update(obj interface{}, db *sqlx.DB) (int64, error) {
    stmts := reflectStatements(obj, true)
    var str string = "UPDATE " + getTable(obj) + " SET " + stmts
    log.Println(str)
    result, err := db.Exec(str)
    nbRows, _ := result.RowsAffected()
    return nbRows, err
}

func Save(obj interface{}, db *sqlx.DB) (int64, error) {
    // Maybe you don't know, but in MySQL at least, you can use the SET syntax like Update for an Insert
    // I've to check if it's compatible with others sql dbs
    stmts := reflectStatements(obj, false)
    var str string = "INSERT INTO " + getTable(obj) + " SET " + stmts
    log.Println(str)
    result, err := db.Exec(str)
    lastId, _ := result.LastInsertId()
    return lastId, err
}
