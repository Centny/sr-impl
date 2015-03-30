package test

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"org.cny.sr/conf"
	"org.cny.sr/dbcon"
)

const TDbCon string = "cny:123@loc.srv:27017/cny"

var Cfg util.Fcfg = util.Fcfg{}

func init() {
	conf.Cfg = Cfg
	ss, err := mgo.Dial(TDbCon)
	if err != nil {
		fmt.Errorf("connection to DB err:%v", err.Error())
		return
	}
	dbcon.Db_ = ss.DB("cny")
}
