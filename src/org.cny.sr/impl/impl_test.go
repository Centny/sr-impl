package impl

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/sr"
	"github.com/Centny/gwf/sr/pb"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"org.cny.sr/dbcon"
	_ "org.cny.sr/test"
	"testing"
)

func TestImpl(t *testing.T) {
	ts := httptest.NewServer(func(hs *routing.HTTPSession) routing.HResult {
		run_t(hs, t)
		return routing.HRES_RETURN
	})
	ts.G("")
}
func run_t(hs *routing.HTTPSession, t *testing.T) {
	dbcon.Db().C("sr").RemoveAll(nil)
	sq := &SrQH{}
	action, name, time, typ := "aa", "bb", util.Now(), int32(101)
	evn := &pb.Evn{}
	evn.Action = &action
	evn.Name = &name
	evn.Time = &time
	evn.Type = &typ
	err := sq.Proc(nil, &sr.SRH_Q_I{
		Sp:   "ss01",
		Aid:  "org.cny",
		Ver:  "0.0.1",
		Dev:  "a",
		Rel:  "rel",
		Evs:  []*pb.Evn{evn},
		Time: util.Now(),
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = sq.Proc(nil, &sr.SRH_Q_I{
		Sp:   "ss02",
		Aid:  "org.cny",
		Ver:  "0.0.1",
		Dev:  "a",
		Rel:  "rel",
		Evs:  []*pb.Evn{evn},
		Time: util.Now(),
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = sq.Proc(nil, &sr.SRH_Q_I{
		Sp:   "ss03",
		Aid:  "org.cny",
		Ver:  "0.0.1",
		Dev:  "a",
		Rel:  "rel",
		Evs:  []*pb.Evn{evn},
		Time: util.Now(),
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	hs.SetVal("rel", "rel")
	v, err := sq.ListSr(nil, hs, "org.cny", "0.0.1", "", "", 0, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}
	vi := v.([]sr.SRH_Q_I)
	fmt.Println(v, "--->1")
	v, err = sq.ListSr(nil, hs, "org.cny", "0.0.1", vi[0].Id.(bson.ObjectId).Hex(), "a", 0, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	vi = v.([]sr.SRH_Q_I)
	fmt.Println(v, "--->2")
	v, err = sq.ListSr(nil, hs, "org.cny", "0.0.1", vi[1].Id.(bson.ObjectId).Hex(), "", 0, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}
	vi = v.([]sr.SRH_Q_I)
	fmt.Println(v)
	v, err = sq.ListSr(nil, hs, "org.cny", "0.0.1", vi[0].Id.(bson.ObjectId).Hex(), "", 0, 0)
	if err == nil {
		t.Error("not error")
		return
	} else if err != util.NOT_FOUND {
		t.Error(err.Error())
	}
	fmt.Println(v)
	//
	fmt.Println("------->")
	v, err = sq.ListPkg(nil, hs, "", "")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	fmt.Println(sq.ListPkg(nil, hs, "org.cny", "a"))
	//
	sq.ListSr(nil, hs, "org.cny", "0.0.1", "5518e0790cbd510cd7000", "", 0, 0)
	sq.Args(nil, nil, "", "", "", "", "")
	NewSr("/tmp")
	//
	hs.SetVal("type", "ssss->")
	sq.ListSr(nil, hs, "org.cny", "0.0.1", "", "", 0, 0)
	hs.SetVal("type", typ)
	hs.SetVal("action", "ssss->")
	sq.ListSr(nil, hs, "org.cny", "0.0.1", "", "", 0, 0)
	hs.SetVal("action", action)
	hs.SetVal("name", name)
	hs.SetVal("type", typ)
	sq.ListSr(nil, hs, "org.cny", "0.0.1", "", "", 0, 0)
	sq.ListSr(nil, hs, "org.cny", "0.0.1", "", "", 0, 1)
}
