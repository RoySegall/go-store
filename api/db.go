package api

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"reflect"
)

func GetSession() (*r.Session) {
	settings := GetSettings()
	session, err := r.Connect(r.ConnectOpts{
		Address: settings.RethinkDB.Address,
		Database: settings.RethinkDB.Database,
	})

	if err != nil {
		panic(err)
	}

	return session
}

// Inserting an object to the DB.
func Insert(table string, arg interface{}) (string) {
	res, err := r.Table(table).Insert(arg).RunWrite(GetSession())

	if err != nil {
		panic(err)
	}

	return res.GeneratedKeys[0]
}

func Get(table string, id string) (*r.Cursor) {
	res, err := r.Table(table).Filter(map[string]interface{} {
		"id": id,
	}).Run(GetSession())

	if err != nil {
		panic(err)
	}

	return res
}

func GetAll(table string) (*r.Cursor) {
	res, err := r.Table(table).OrderBy("id").Run(GetSession())

	if err != nil {
		panic(err)
	}

	return res
}

func Update(table string, object interface{}) (*r.Cursor) {
	res, err := r.Table(table).Get(reflect.Indirect(reflect.ValueOf(object))).Update(object).Run(GetSession())

	if err != nil {
		panic(err)
	}

	return res
}

func Delete(table string, id string) (*r.Cursor) {
	res, err := r.Table(table).Filter(map[string]interface{} {
		"id": id,
	}).Delete().Run(GetSession())

	if err != nil {
		panic(err)
	}

	return res
}

// Create a DB.
func DbCreate(db string) {
	err := r.DBCreate(db).Exec(GetSession())

	if err != nil {
		panic(err)
	}
}

// Checking if a DB exists.
func DbDrop(db string) {
	err := r.DBDrop(db).Exec(GetSession())

	if err != nil {
		panic(err)
	}
}

// Create a table.
func TableCreate(table string) {
	err := r.TableCreate(table).Exec(GetSession())

	if err != nil {
		panic(err)
	}
}

// Dropping a table
func TableDrop(table string) {
	err := r.TableDrop(table).Exec(GetSession())

	if err != nil {
		panic(err)
	}
}
