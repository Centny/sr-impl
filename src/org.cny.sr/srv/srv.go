//Package provide server function.
//Author:Centny
package srv

import (
	"fmt"
	"github.com/Centny/gwf/im"
	"github.com/Centny/gwf/log"
	"gopkg.in/mgo.v2"
	"net/http"
	"org.cny.sr/conf"
	"org.cny.sr/dbcon"
	"runtime"
	"sync"
)

var lock sync.WaitGroup
var s_running bool
var l *im.Listener

func run(args []string) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	defer StopSrv()
	cfile := "conf/sr.properties"
	if len(args) > 1 {
		cfile = args[1]
	}
	fmt.Println("Using config file:", cfile)
	err := conf.Cfg.InitWithFilePath(cfile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//
	log.I("Config:\n%v", conf.Cfg.Show())
	//test connect
	if len(conf.SrDbConn()) < 1 {
		fmt.Println("SR_DB_CONN is not exist in config")
		return
	}
	ss, err := mgo.Dial(conf.SrDbConn())
	if err != nil {
		fmt.Errorf("connection to DB err:%v", err.Error())
		return
	}
	dbcon.Con_ = ss
	dbcon.DbName_ = conf.SrDbName()
	//
	mux := http.NewServeMux()
	HSrvMux(mux, "", "www")
	log.D("running web server on %s", conf.ListenAddr())
	s := http.Server{Addr: conf.ListenAddr(), Handler: mux}
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

//run the server.
func RunSrv(args []string) {
	s_running = true
	lock.Add(1)
	go run(args)
	lock.Wait()
	s_running = false
}

//stop the server.
func StopSrv() {
	if s_running {
		lock.Done()
	}
}
