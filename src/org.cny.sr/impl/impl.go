package impl

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/sr"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"org.cny.sr/dbcon"
)

type SrQH struct {
}

func (srq *SrQH) Args(s *sr.SRH_Q, hs *routing.HTTPSession, aid, ver, dev, sp, sf string) (util.Map, error) {
	return nil, nil
}
func (srq *SrQH) Proc(s *sr.SRH_Q, i *sr.SRH_Q_I) error {
	i.Id = bson.NewObjectId()
	// log.D("adding SRH_Q_I:%v", i)
	return dbcon.Db().C("sr").Insert(i)
}
func (srq *SrQH) ListSr(s *sr.SRH_Q, hs *routing.HTTPSession, aid, ver, prev, dev string, from, all int64) (interface{}, error) {
	var action, name string
	var typ int64 = math.MinInt64
	err := hs.ValidCheckVal(`
		action,O|S,L:0;
		name,O|S,L:0;
		type,O|I,R:0;
		`, &action, &name, &typ)
	if err != nil {
		return nil, err
	}
	qa := bson.M{
		"aid": aid,
		"ver": ver,
		"time": bson.M{
			"$gt": from,
		},
	}
	if len(prev) > 0 {
		if !bson.IsObjectIdHex(prev) {
			return nil, util.Err("invalid prev id(%v)", prev)
		}
		qa["_id"] = bson.M{
			"$gt": bson.ObjectIdHex(prev),
		}
	}
	if len(dev) > 0 {
		qa["dev"] = dev
	}
	tmatch := bson.M{}
	if len(action) > 0 {
		tmatch["evs.action"] = action
	}
	if len(name) > 0 {
		tmatch["evs.name"] = name
	}
	if typ > math.MinInt64 {
		tmatch["evs.type"] = typ
	}
	var res_l []sr.SRH_Q_I = []sr.SRH_Q_I{}
	var res_i sr.SRH_Q_I
	log.I("list SR by all(%v) from(%v) %v->%v", all, from, util.S2Json(qa), util.S2Json(tmatch))
	// var query *mgo.Query
	if len(tmatch) > 0 {
		pipe := []bson.M{
			bson.M{
				"$match": qa,
			}, bson.M{
				"$unwind": "$evs",
			}, bson.M{
				"$match": tmatch,
			}, bson.M{
				"$group": bson.M{
					"_id": bson.M{
						"_id":  "$_id",
						"sp":   "$sp",
						"aid":  "$aid",
						"ver":  "$ver",
						"dev":  "$dev",
						"time": "$time",
					},
					"evs": bson.M{
						"$push": "$evs",
					},
				},
			}, bson.M{
				"$project": bson.M{
					"_id":  "$_id._id",
					"sp":   "$_id.sp",
					"aid":  "$_id.aid",
					"ver":  "$_id.ver",
					"dev":  "$_id.dev",
					"time": "$_id.time",
					"evs":  "$evs",
				},
			},
		}
		if all > 0 {
			err = dbcon.Db().C("sr").Pipe(pipe).All(&res_l)
			return res_l, err
		} else {
			err = dbcon.Db().C("sr").Pipe(pipe).One(&res_i)
			if err == mgo.ErrNotFound {
				err = util.NOT_FOUND
			}
			return []sr.SRH_Q_I{res_i}, err
		}
	} else {
		if all > 0 {
			err = dbcon.Db().C("sr").Find(qa).All(&res_l)
			return res_l, err
		} else {
			err = dbcon.Db().C("sr").Find(qa).One(&res_i)
			if err == mgo.ErrNotFound {
				err = util.NOT_FOUND
			}
			return []sr.SRH_Q_I{res_i}, err
		}
	}
}

func NewSr(r string) (*sr.SR, *sr.SRH_Q) {
	return sr.NewSR3(r, &SrQH{})
}
