package srv

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/doc"
	"github.com/Centny/gwf/routing/filter"
	"net/http"
	"org.cny.sr/impl"
	"org.cny.sr/mr"
	"regexp"
	"runtime"
)

func HSrvMux(smux *http.ServeMux, pre string, www string) {
	mux := routing.NewSessionMux2(pre)
	// mux.ShowLog = true
	cors := filter.NewCORS()
	cors.AddSite("*")
	mux.HFilter("^/.*$", cors)
	sr, srq := impl.NewSr(www + "/sdata")
	srq.Run(runtime.NumCPU() - 1)
	log.D("register sr...")
	mux.H("^/sr(\\?.*)?$", sr)
	log.D("register mr...")
	mux.H("^/mr(/.*)?$", mr.NewMR("/mr"))
	//
	//
	dv := doc.NewDocViewer()
	dv.Excs = []*regexp.Regexp{
		regexp.MustCompile("^\\^/\\.\\*\\$$"),
		regexp.MustCompile("^\\^/api/doc\\.\\*\\$$"),
	}
	mux.H("^/api/doc.*$", dv)
	//
	mux.Handler("^/.*$", http.FileServer(http.Dir(www)))
	//
	smux.Handle("/", mux)
}
