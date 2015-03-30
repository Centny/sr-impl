package dbcon

import (
	"gopkg.in/mgo.v2"
)

var Db_ *mgo.Database

func Db() *mgo.Database {
	if Db_ == nil {
		panic("database not initial")
	}
	return Db_
}
