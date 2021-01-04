package broker

import (
	"log"

	"github.com/masagatech/myinvest/module"

	"github.com/masagatech/myinvest/db"
	"github.com/masagatech/myinvest/util"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo/v4"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	ur := r.Group("/broker")
	ur.GET("/get/:uid", get)
	ur.POST("/login", login)

}

type broker struct {
	Code        string `bson:"code" json:"code"`
	Name        string `bson:"name" json:"name"`
	ShortInfo   string `bson:"shortinfo" json:"shortinfo"`
	ApiURL      string `bson:"apiurl" json:"apiurl"`
	ApiKey      string `bson:"apikey,omitempty" json:"apikey,omitempty"`
	Logo        string `bson:"logo,omitempty" json:"logo,omitempty"`
	Description string `bson:"desc" json:"desc"`
	Active      bool   `bson:"active" json:"active"`
}

func get(c echo.Context) error {

	_userid := c.Param("uid")
	//var _broker []broker
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	_col, _ := db.GetCollection(c, db.Broker)

	_query := []bson.M{{
		"$lookup": bson.M{ // lookup the documents table here
			"from": db.Thirdpartytokens,
			"let":  bson.M{"code": "$code"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$and": []bson.M{
							{
								"$eq": []interface{}{"$userid", bson.ObjectIdHex(_userid)},
							}, {
								"$eq": []string{"$$code", "$code"},
							},
						},
					},
				},
				},
			},
			"as": "b",
		}},
		bson.M{
			"$match": bson.M{
				"active": true,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"preserveNullAndEmptyArrays": true,
				"path":                       "$b",
			},
		}, bson.M{
			"$project": bson.M{
				"_id":       0,
				"userid":    "$b.userid",
				"active":    1,
				"code":      1,
				"name":      1,
				"shortinfo": 1,
				"logo":      1,
				"desc":      1,
			},
		},
	}

	var brokerlist []interface{}
	pipe := _col.Pipe(_query)
	err := pipe.All(&brokerlist)

	if err != nil {
		log.Fatal(err)
	}

	return util.Response(c, 1, "", "", &brokerlist)

}

func login(c echo.Context) error {

	params := struct {
		Broker string `bson:"broker" json:"broker"`
		Userid string `bson:"uid" json:"uid"`
	}{}

	if err := c.Bind(&params); err != nil {
		return util.Response(c, 0, "001", err.Error(), bson.M{})
	}

	login_url := getZerodhaLoginUrl(c)
	login_url += "&redirect_params=uid%3D" + params.Userid

	return util.Response(c, 1, "", "", login_url)

}

func getZerodhaLoginUrl(c echo.Context) string {
	kc := module.GetZerodhaInstance(c)
	return kc.GetLoginURL()

}
