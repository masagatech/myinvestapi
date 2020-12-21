package order

import (
	"encoding/csv"
	"log"
	"strconv"

	"github.com/masagatech/myinvest/db"
	"github.com/masagatech/myinvest/module"
	"github.com/masagatech/myinvest/util"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	ur := r.Group("/orders")
	ur.POST("/create", create)
}

type Instrument struct {
	Instrument_Token string  `bson:"instrument_token"`
	Exchange_Token   string  `bson:"exchange_token"`
	Tradingsymbol    string  `bson:"tradingsymbol"`
	Name             string  `bson:"name"`
	Last_Price       float64 `bson:"last_price"`
	Expiry           string  `bson:"expiry"`
	Strike           float64 `bson:"strike"`
	Tick_Size        float64 `bson:"tick_size"`
	Lot_Size         int     `bson:"lot_size"`
	Instrument_type  string  `bson:"instrument_type"`
	Segment          string  `bson:"segment"`
	Exchange         string  `bson:"exchange"`
}

type GTTOrder struct {
	Exchange         string    `json:"exchange"`
	Tradingsymbol    string    `json:"tradingsymbol"`
	Trigger_values   []float64 `json:"trigger_values"`
	Last_price       float64   `json:"last_price"`
	Transaction_type string    `json:"transaction_type"`
	Quantity         int       `json:"quantity"`
	Order_type       string    `json:"order_type"`
	Product          string    `json:"product"`
	Price            float64   `json:"price"`
	Uid              string    `json:"uid"`
}

const (
	apiKey    string = "my_api_key"
	apiSecret string = "my_api_secret"
)

func create(c echo.Context) error {

	orderReq := GTTOrder{}

	if err := c.Bind(&orderReq); err != nil {
		fmt.Println(err.Error())
	}

	kc := module.GetZerodhaInstance(c)
	// kc := kiteconnect.New(apiKey)
	// kc.GetLoginURL()
	ord := kiteconnect.GTTParams{
		Tradingsymbol:   orderReq.Tradingsymbol,
		Exchange:        orderReq.Exchange,
		LastPrice:       orderReq.Last_price,
		TransactionType: orderReq.Order_type,
		Trigger: &kiteconnect.GTTSingleLegTrigger{
			TriggerParams: kiteconnect.TriggerParams{
				TriggerValue: 1,
				Quantity:     1,
				LimitPrice:   1,
			},
		},
	}

	gttResp, err := kc.PlaceGTT(ord)
	if err != nil {
		log.Fatalf("error placing gtt: %v", err)
	}

	log.Println("placed GTT trigger_id = ", gttResp.TriggerID)

	log.Println("Fetching details of placed GTT...")

	// order, err := kc.GetGTT(gttResp.TriggerID)
	// if err != nil {
	// 	log.Fatalf("Error getting GTTs: %v", err)
	// }

	exchange := c.Param("exchange")

	response, err := http.Get(util.Config.Zerodha.ApiUrl + "/instruments/" + exchange)

	if err != nil {
		fmt.Print(err.Error())
		return util.Response(c, 1, "001", err.Error(), bson.M{})
	}

	// responseData, err := ioutil.ReadAll(response.Body)
	r := csv.NewReader(response.Body)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		return util.Response(c, 1, "002", err.Error(), bson.M{})
	}
	_col, _ := db.GetCollection(c, db.Instruments)
	_blk := _col.Bulk()
	//str := "2014-11-12"
	var _instruments []interface{}
	//fmt.Println(len(records[1:]))
	_total_records := len(records) - 1
	fmt.Println(_total_records)
	for _, rec := range records[1:] {
		var _instrument Instrument
		_instrument.Instrument_Token = rec[0]
		_instrument.Exchange_Token = rec[1]
		_instrument.Tradingsymbol = rec[2]
		_instrument.Name = rec[3]
		_instrument.Last_Price, _ = strconv.ParseFloat(rec[4], 64)
		_instrument.Expiry = rec[5]
		_instrument.Strike, _ = strconv.ParseFloat(rec[6], 64)
		_instrument.Tick_Size, _ = strconv.ParseFloat(rec[7], 64)
		_instrument.Lot_Size, _ = strconv.Atoi(rec[8])
		_instrument.Instrument_type = rec[9]
		_instrument.Segment = rec[10]
		_instrument.Exchange = rec[11]

		// _instrument.EmployeeName = rec[1]
		// emp.EmployeeSalary, = strconv.Atoi(rec[2])
		// emp.EmployeeAge, _ = strconv.Atoi(rec[3])
		// emp.ProfileImage = rec[4]

		_selector := bson.M{"instrument_token": _instrument.Instrument_Token}
		//_p := Pair{_selector, _instrument}
		_instruments = append(_instruments, _selector)
		_instruments = append(_instruments, _instrument)
		// _, err = _col.Upsert(_selector, _instrument)
		// if err != nil {
		// 	log.Fatal(err)
		// }

	}
	//_selector := bson.M{"instrument_token": _instrument.Instrument_Token}
	_blk.Unordered()
	_blk.Upsert(_instruments...)
	res, err := _blk.Run()
	fmt.Println(res)
	fmt.Println(err)

	return util.Response(c, 1, "", fmt.Sprintf("Total %d records synced", _total_records), bson.M{})

}
