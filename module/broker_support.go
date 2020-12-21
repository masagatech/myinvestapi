package module

import (
	"github.com/masagatech/myinvest/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Brokers_init(_db *mgo.Session) {

	_col, _ := db.GetDBCollection(_db, db.Broker)

	// zerodha
	_col.Upsert(bson.M{
		"code": "zerodha",
	}, bson.M{
		"code":      "zerodha",
		"name":      "Kite Connect",
		"shortinfo": "zerodha broker",
		"apiurl":    "",
		"apikey":    "",
		"logo":      "data:image/png;base64,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",
		"desc":      "India's biggest stock broker offering the lowest, cheapest brokerage rates for futures and options, commodity trading, equity and mutual funds.",
		"active":    true,
	})

}
