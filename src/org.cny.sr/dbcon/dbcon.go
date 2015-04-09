package dbcon

import (
	"gopkg.in/mgo.v2"
)

var Con_ *mgo.Session
var DbName_ string

func Db() *mgo.Database {
	if Con_ == nil {
		panic("database not initial")
	}
	if err := Con_.Ping(); err != nil {
		Con_.Refresh()
	}
	return Con_.DB(DbName_)
}
