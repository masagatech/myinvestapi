package strategy

import (
	"fmt"
	"log"

	"github.com/masagatech/myinvest/api/integration"
	kiteconnect "github.com/zerodhatech/gokiteconnect"

	"github.com/masagatech/myinvest/util"

	"github.com/masagatech/myinvest/db"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo/v4"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	ur := r.Group("/strategy")
	ur.POST("/upsert", create)
	ur.GET("/get", get)
	ur.GET("/edit/:id", edit)
	ur.DELETE("/:id", delete)
	ur.GET("/start/:id", start)

}

type Strategy struct {
	Id                    string  `bson:"-" json:"_id"`
	Name                  string  `bson:"name"`
	Type                  string  `bson:"type"`
	Gtt_type              string  `bson:"gtt_type"`
	Expiry                string  `bson:"expiry"`
	Instrument_id         int     `bson:"instrument_id"`
	Instrument_name       string  `bson:"instrument_name"`
	Qty                   float64 `bson:"qty"`
	Order_price           float64 `bson:"order_price"`
	Trigger_order_price   float64 `bson:"trigger_order_price"`
	Trigger_point         int     `bson:"trigger_point"`
	Stop_loss_price       float64 `bson:"stop_loss_price"`
	Stop_loss_point       int     `bson:"stop_loss_point"`
	Stop_loss_order_level float64 `bson:"stop_loss_order_level"`
	User_created          string  `bson:"user_created"`
	IsStart               bool    `bson:"isstarted"`
	Tradingsymbol         string  `bson:"tradingsymbol"`
	Exchange              string  `bson:"exchange"`
	TriggerId             int     `bson:"triggerid"`
}

func create(c echo.Context) error {

	_col, _ := db.GetCollection(c, db.Strategy)
	_strategy := Strategy{}

	_id := bson.NewObjectId()
	err := c.Bind(&_strategy)

	if _strategy.Id != "" {
		_id = bson.ObjectIdHex(_strategy.Id)
	}
	fmt.Println(_strategy.Id)
	if err != nil {
		return util.Response(c, 0, "", err.Error(), bson.M{})
	}
	d, err1 := _col.Upsert(bson.M{"_id": _id}, &_strategy)
	if err1 != nil {
		return util.Response(c, 0, "", err1.Error(), bson.M{})
	}
	iu := 0
	iumsg := "Inserted successfully."
	if d.Updated > 0 {
		iu = 1
		iumsg = "Updated successfully."
	}

	return util.Response(c, 1, "", iumsg, bson.M{
		"iu": iu,
		"id": d.UpsertedId,
	})

}

func get(c echo.Context) error {
	_col, _ := db.GetCollection(c, db.Strategy)
	var _strategy []interface{}
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	err := _col.Find(nil).Select(bson.M{"_id": 1,
		"name":                  1,
		"gtt_type":              1,
		"instrument_name":       1,
		"qty":                   1,
		"order_price":           1,
		"isstarted":             1,
		"expiry":                1,
		"stop_loss_price":       1,
		"stop_loss_point":       1,
		"trigger_point":         1,
		"stop_loss_order_level": 1,
		"tradesymbol":           1,
		"exchange":              1,
	}).All(&_strategy)

	// _, err := _db.QueryOne(&_user, sql)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err)
		return util.Response(c, 0, "", err.Error(), bson.M{})

	}

	return util.Response(c, 1, "", "", &_strategy)
}

func edit(c echo.Context) error {
	_col, _ := db.GetCollection(c, db.Strategy)
	id := c.Param("id")
	var _strategy interface{}
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	_col.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1,
		"name":                  1,
		"gtt_type":              1,
		"instrument_name":       1,
		"instrument_id":         1,
		"qty":                   1,
		"order_price":           1,
		"isstarted":             1,
		"expiry":                1,
		"stop_loss_price":       1,
		"stop_loss_point":       1,
		"trigger_point":         1,
		"stop_loss_order_level": 1,
		"tradesymbol":           1,
		"exchange":              1,
	}).One(&_strategy)

	// _, err := _db.QueryOne(&_user, sql)

	return util.Response(c, 1, "", "", &_strategy)
}

func delete(c echo.Context) error {
	_col, _ := db.GetCollection(c, db.Strategy)
	id := c.Param("id")
	var _strategy interface{}
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	_col.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1,
		"name":                  1,
		"gtt_type":              1,
		"instrument_name":       1,
		"instrument_id":         1,
		"qty":                   1,
		"order_price":           1,
		"isstarted":             1,
		"expiry":                1,
		"stop_loss_price":       1,
		"stop_loss_point":       1,
		"trigger_point":         1,
		"stop_loss_order_level": 1,
		"tradesymbol":           1,
		"exchange":              1,
	}).One(&_strategy)

	// _, err := _db.QueryOne(&_user, sql)

	return util.Response(c, 1, "", "", &_strategy)
}

func start(c echo.Context) error {

	_col, _ := db.GetCollection(c, db.Strategy)
	id := c.Param("id")
	var _strategy Strategy
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	_col.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1,
		"name":                  1,
		"gtt_type":              1,
		"instrument_name":       1,
		"instrument_id":         1,
		"qty":                   1,
		"order_price":           1,
		"isstarted":             1,
		"expiry":                1,
		"stop_loss_price":       1,
		"stop_loss_point":       1,
		"trigger_point":         1,
		"stop_loss_order_level": 1,
		"tradesymbol":           1,
		"exchange":              1,
	}).One(&_strategy)

	// kc := module.GetZerodhaInstance(c)

	// _col_third_party, _ := db.GetCollection(c, db.Thirdpartytokens)

	return util.Response(c, 1, "", "", &_strategy)

}

func placeGTT(_strategy Strategy, c echo.Context) int {
	_session_token := util.GetTokenData(c)
	kc := integration.GetSesssion(c, _session_token.Id)
	var transactionType = ""
	var triggerValue = 0.0
	var limitPrice = 0.0

	if _strategy.Gtt_type == "BUY" {
		if _strategy.Type == "BUY" {
			transactionType = kiteconnect.TransactionTypeBuy
			triggerValue = _strategy.Order_price
			limitPrice = _strategy.Stop_loss_price
		} else {
			transactionType = kiteconnect.TransactionTypeSell
			triggerValue = _strategy.Order_price
			limitPrice = _strategy.Stop_loss_order_level
		}
	} else if _strategy.Gtt_type == "SELL" {
		if _strategy.Type == "SELL" {
			transactionType = kiteconnect.TransactionTypeBuy
			triggerValue = _strategy.Order_price
			limitPrice = _strategy.Stop_loss_price
		} else {
			transactionType = kiteconnect.TransactionTypeSell
			triggerValue = _strategy.Order_price
			limitPrice = _strategy.Stop_loss_order_level
		}

	}

	gttResp, err := kc.PlaceGTT(kiteconnect.GTTParams{
		Tradingsymbol:   _strategy.Tradingsymbol,
		Exchange:        _strategy.Exchange,
		LastPrice:       800,
		TransactionType: transactionType,
		Trigger: &kiteconnect.GTTSingleLegTrigger{
			TriggerParams: kiteconnect.TriggerParams{
				TriggerValue: triggerValue,
				Quantity:     _strategy.Qty,
				LimitPrice:   limitPrice,
			},
		},
	})
	if err != nil {
		log.Fatalf("error placing gtt: %v", err)
	}

	log.Println("placed GTT trigger_id = ", gttResp.TriggerID)
	return gttResp.TriggerID

}
