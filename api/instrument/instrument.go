package instrument

import (
	"encoding/csv"
	"log"
	"strconv"

	"github.com/masagatech/myinvest/db"
	"github.com/masagatech/myinvest/util"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	ur := r.Group("/instruments")
	ur.GET("/sync/:exchange", sync)
	ur.GET("/autocomplete/:query", autocomplete)

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

type Pair struct {
	a, b interface{}
}

func sync(c echo.Context) error {

	exchange := c.Param("exchange")

	fmt.Println(util.Config.Zerodha.ApiUrl + "/instruments/" + exchange)

	response, err := http.Get(util.Config.Zerodha.ApiUrl + "/instruments/" + exchange)

	if err != nil {
		fmt.Print(err.Error())
		return util.Response(c, 1, "001", err.Error(), bson.M{})
	}

	// responseData, err := ioutil.ReadAll(response.Body)
	r := csv.NewReader(response.Body)
	fmt.Println(response.Body)
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

func autocomplete(c echo.Context) error {

	query := c.Param("query")

	var _instr []interface{}
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	_col, _ := db.GetCollection(c, db.Instruments)
	err := _col.Find(bson.M{"name": bson.RegEx{query, "i"}}).Select(bson.M{"_id": 0, "name": 1,
		"instrument_token": 1, "instrument_type": 1,
		"expiry": 1, "exchange": 1,
		"tradingsymbol": 1,
	}).Limit(8).All(&_instr)
	if err != nil {
		log.Fatal(err)
	}
	// pipeLine := []bson.M{
	// 	bson.M{"$project": bson.M{"name": 1, "instrument_token": 1}},   // output address, city, st, notecount
	// 	bson.M{"$match": bson.M{"name": "WINSOME TEXTILE INDUSTRIES"}}, // keep docs with more than 2 notes
	// 	bson.M{"$sort": bson.D{{"name", 1}}},                           // sort results by state, city - see note above
	// }

	// iter := _col.Pipe(pipeLine).Iter()

	// defer iter.Close()

	// for iter.Next(&result) {
	// 	log.Printf("%+v", result)
	// }

	return util.Response(c, 1, "", "", &_instr)

}
