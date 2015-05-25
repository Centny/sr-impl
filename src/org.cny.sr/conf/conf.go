package conf

import (
	"github.com/Centny/gwf/util"
)

var Cfg *util.Fcfg = util.NewFcfg3()

func ListenAddr() string {
	return Cfg.Val("LISTEN_ADDR")
}

func SrDbConn() string {
	return Cfg.Val("SR_DB_CONN")
}
func SrDbName() string {
	return Cfg.Val("SR_DB_NAME")
}
