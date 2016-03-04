# go-save-model-sqlx

## Introduction
This is an experimental and test sample of code, I was just wondering how to produce a generic Save method for my models while using sqlx package. This is perfectly suitable for database/sql package, just to change few parameters..

DO NOT USE IN PRODUCTION

I've tried to handle common sql types, but there is probably things that won't work, any help is welcome :P

The goal here, is to be able to do something like this :
```golang
u := User{}
... // whatever how you fill your object
u.Save()
```
I didn't wanted to produce a separate object, struct of control or anything. I could've done a UserMapper object which is charged to save my primitive model. Or any other pattern. My goal was to be close of common ORM you can find in other languages.
This is not perfect, but it's working :)

## Install
`go get -u github.com/snickers54/go-save-model-sqlx`
`import github.com/snickers54/go-save-model-sqlx/models`
## HOWTO
The only file you want to include in your project, is `models/model.go`. The others are just here to create a viable example of use.
I've based my runtime reflection on a common tag "db" used by others. I could've done a go generate and based my reflection on compile time stuff. But didn't wanted to.

To be compatible, your models need to implement the interface Model which is quite simple, because you just have to implement 2 methods : `Save()` and `Table()`..

The tag `db:"fieldName"` is used to determine what the field name is .. I wanted it to be explicit, I can't just base it on the name of your field in the struct ..

By default, you can build, but it won't run, the `models/mysql.go` file is trying to connect on a mysql database which probably doesn't exists ..

## Optional
### AutoIncr struct
This struct is for my personal use, really often needed in my models, it's only about a auto-incremental primary key and two fields representing creation and update dates..
I let it in `models/model.go` because I think it's interesting and needed in my example to show how it's handled by the reflectStatements function.

### Driver struct
This struct is by default embbed into the AutoIncr struct.
This is an empty struct, no fields at all, only a method to get an instance of sql.DB or sqlx.DB. It's needed in my example to show that the reflection is not gonna use it because it's an empty struct. But it's also really useful, because I'm able to do something like that :
```golang
u := User{}
u.db().Exec("...")
```
I didn't export the symbol, because you should not use this in anything else than the `models` package.
