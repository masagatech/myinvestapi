package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/masagatech/myinvest/db"
	"github.com/masagatech/myinvest/util"
	"gopkg.in/mgo.v2/bson"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	r.POST("/auth", auth)
}

func auth(c echo.Context) (err error) {
	_config := util.GetConfig(c)
	var _user struct {
		Id       bson.ObjectId `json:"id" bson:"_id"`
		Name     string        `json:"name" bson:"name"`
		Email    string        `json:"email" bson:"email"`
		Password string        `json:"-" bson:"password,omitempty"`
		Active   bool          `json:"-" bson:",omitempty"`
		Token    string        `json:"token" bson:"-"`
	}
	_params := echo.Map{}
	if err := c.Bind(&_params); err != nil {
		return err
	}

	email := _params["email"].(string)
	password := _params["password"].(string)
	_col, _ := db.GetCollection(c, db.Users)
	err1 := _col.Find(bson.M{
		"email": email,
	}).Select(bson.M{
		"_id":      1,
		"name":     1,
		"email":    1,
		"password": 1,
		"active":   1,
	}).One(&_user)

	if err1 != nil {
		fmt.Println(err)
	}

	resultKey := 1
	errorCode := ""
	message := ""

	if _user.Email == "" {
		resultKey = 0
		errorCode = "001"
		message = "Invalid Username/Password"
	} else if _user.Active == false {
		resultKey = 0
		errorCode = "003"
		message = "User is blocked. Please contact administrator"
	} else if _user.Password != password {
		resultKey = 0
		errorCode = "002"
		message = "Invalid Username/Password"
	}

	if resultKey == 1 {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = _user.Id
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		to, err := token.SignedString([]byte(_config.JWT.Secret))
		if err != nil {
			fmt.Println("Error")
		}
		_user.Token = to
	}

	return util.Response(c, resultKey, errorCode, message, &_user)
}
