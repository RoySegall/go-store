package api

import (
	r "gopkg.in/gorethink/gorethink.v3"
)

type TodoItem struct {
	Id string
}

func getSession() (*r.Session) {
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost",
		Database: "store",
	})

	if err != nil {
		panic(err)
	}

	return session
}

// Inserting an object to the DB.
func Insert(table string, arg interface{}) (string) {
	res, err := r.Table(table).Insert(arg).RunWrite(getSession())

	if err != nil {
		panic(err)
	}

	return res.GeneratedKeys[0]
}

func Get(table string, id string) (*r.Cursor) {
	res, err := r.Table(table).Filter(map[string]interface{}{
		"id": id,
	}).Run(getSession())

	if err != nil {
		panic(err)
	}

	return res
}

func GetAll(table string) (*r.Cursor) {
	res, err := r.Table(table).Run(getSession())

	if err != nil {
		panic(err)
	}

	return res
}

func Delete(table string, id string) (*r.Cursor) {
	res, err := r.Table(table).Filter(map[string]interface{} {
		"id": id,
	}).Delete().Run(getSession())

	if err != nil {
		panic(err)
	}

	return res
}

// Create a DB.
func DbCreate(db string) {
	err := r.DBCreate(db).Exec(getSession())

	if err != nil {
		panic(err)
	}
}

// Create a table.
func TableCreate(table string) {
	err := r.TableCreate(table).Exec(getSession())

	if err != nil {
		panic(err)
	}
}
