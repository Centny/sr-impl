package mr

import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"org.cny.sr/dbcon"
	"strconv"
	"strings"
)

type MR_S struct {
}

type MR struct {
	Pre string
	Did string //the default id
}

func NewMR(pre string) *MR {
	if len(pre) < 1 {
		pre = "/"
	}
	return &MR{
		Pre: pre,
		Did: "_default",
	}
}
func (m *MR) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	path := strings.TrimPrefix(hs.R.URL.Path, m.Pre)
	path = strings.Trim(path, "/ \t")
	path = strings.Replace(path, "/", ".", -1)
	var typ, data, exec, id string = "S", "", "G", m.Did
	err := hs.ValidF(`
		id,O|S,L:0;
		type,O|S,O:I~F~S~J;
		data,O|S,L:0;
		exec,O|S,O:S~G~I~D;
		`, &id, &typ, &data, &exec)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	var val interface{}
	switch exec {
	case "S":
		err = m.Set(id, path, typ, data)
		val = "Ok"
	case "I":
		err = m.Inc(id, path, typ, data)
		val = "OK"
	case "D":
		err = m.Del(id, path)
		val = "OK"
	default:
		val, err = m.Get(id, path)
	}
	if err == nil {
		return hs.MsgRes(val)
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}
func (m *MR) Set(id, path, typ, data string) error {
	var val interface{}
	var err error
	switch typ {
	case "I":
		val, err = strconv.ParseInt(data, 10, 64)
	case "F":
		val, err = strconv.ParseFloat(data, 64)
	case "S":
		val = data
	default:
		val, err = util.Json2Map(data)
	}
	if err != nil {
		return err
	}
	_, err = dbcon.Db().C("mrs").Upsert(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			path: val,
		},
	})
	return err
}
func (m *MR) Get(id, path string) (interface{}, error) {
	var mv util.Map
	if len(path) < 1 {
		return &mv, dbcon.Db().C("mrs").Find(bson.M{"_id": id}).One(&mv)
	} else {
		err := dbcon.Db().C("mrs").Pipe([]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": id,
				},
			},
			bson.M{
				"$project": bson.M{
					"val": "$" + path,
				},
			},
		}).One(&mv)
		return mv.Val("val"), err
	}
}
func (m *MR) Del(id, path string) error {
	if len(path) < 1 {
		return util.Err("please specified attr path")
	} else {
		return dbcon.Db().C("mrs").Update(
			bson.M{
				"_id": id,
			},
			bson.M{
				"$unset": bson.M{
					path: 1,
				},
			},
		)
	}
}
func (m *MR) Inc(id, path, typ, data string) error {
	var val interface{}
	var err error
	switch typ {
	case "I":
		val, err = strconv.ParseInt(data, 10, 64)
	case "F":
		val, err = strconv.ParseFloat(data, 64)
	default:
		err = util.Err("not support Inc for type(%v)", typ)
	}
	if err != nil {
		return err
	}
	return dbcon.Db().C("mrs").Update(
		bson.M{
			"_id": id,
		},
		bson.M{
			"$inc": bson.M{
				path: val,
			},
		},
	)
}
