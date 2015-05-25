package mr

import (
	"fmt"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"org.cny.sr/dbcon"
	_ "org.cny.sr/test"
	"testing"
)

func TestMr(t *testing.T) {
	dbcon.Db().C("mrs").RemoveAll(bson.M{})
	mm := NewMR("")
	ts := httptest.NewServer2(mm)
	mv, _ := ts.G2("/abc/x/a?exec=S&data=1&type=I")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/abc/x/b?exec=S&data=100&type=F")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/abc/x/c?exec=S&data=100&type=S")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/abc/x/d?exec=S&data=%v&type=J", util.S2Json(util.Map{
		"x": 1,
		"y": "2",
		"D": util.Map{
			"a": 1,
			"b": 2,
		},
	}))
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("")
	if mv.IntValP("/data/abc/x/a") != 1 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	if mv.IntValP("/data/abc/x/b") != 100 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	if mv.StrValP("/data/abc/x/c") != "100" {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/abc/x/d")
	if mv.IntValP("/data/x") != 1 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	fmt.Println(mv)
	//
	mv, _ = ts.G2("/xxx/a?exec=S&data=abc")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/xxx/a")
	if mv.StrValP("/data") != "abc" {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/xxx")
	if mv.StrValP("/data/a") != "abc" {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	//
	//test inc
	mv, _ = ts.G2("/abc/x/b?exec=I&data=100&type=F")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/abc/x/a?exec=I&data=100&type=I")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/xxd/sd/a?exec=I&data=100&type=I")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("/xxd/sd/a?exec=I&data=100&type=I")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("")
	if mv.IntValP("/data/abc/x/a") != 101 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	if mv.IntValP("/data/abc/x/b") != 200 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	if mv.IntValP("/data/xxd/sd/a") != 200 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	//test delete
	mv, _ = ts.G2("/abc/x/a?exec=D")
	if mv.IntVal("code") != 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("")
	if val, err := mv.ValP("/data/abc/x/a"); err == nil || val != nil {
		fmt.Println(mv, err, val)
		t.Error("error")
		return
	}
	mv, _ = ts.G2("?exec=D")
	if mv.IntVal("code") == 0 {
		fmt.Println(mv)
		t.Error("error")
		return
	}
	// mv, _ = ts.G2("")
	// fmt.Println(mv)

	//test error
	fmt.Println(ts.G2("/axx/ss"))
	fmt.Println(ts.G2("/axx/ss?exec=sss"))
	fmt.Println(ts.G2("/axx/ss?exec=S&data=xxd&type=I"))
	fmt.Println(ts.G2("/axx/ss?exec=I&data=xxd&type=I"))
	fmt.Println(ts.G2("/axx/ss?exec=I&data=1&type=S"))
}
