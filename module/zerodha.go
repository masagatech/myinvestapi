package module

import (
	"github.com/labstack/echo/v4"
	"github.com/masagatech/myinvest/util"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

func GetZerodhaInstance(c echo.Context) *kiteconnect.Client {
	_config := c.Get("config").(*util.Configuration)
	//util.SessionGet(c, )
	//_col, _ := db.GetCollection(c, db.Thirdpartytokens)

	//var _token_data interface{}
	// query := _col.Find(bson.M{
	// 	"userid": _userid,
	// })

	// if er := query.One(&_token_data); er != nil {
	// 	fmt.Println("error", er.Error())
	// }

	kc := kiteconnect.New(_config.Zerodha.ApiKey)
	kc.SetBaseURI(_config.Zerodha.ApiUrl)
	//kc.SetAccessToken()

	//us, err := kc.GenerateSession("", "")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	return kc
}
