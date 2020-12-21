package db

import (
	"log"

	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2"
)

var g_database = ""

// New creates new database connection to a postgres database
func New(psn string, database string) (*mgo.Session, error) {
	g_database = database
	_db, err := mgo.Dial(psn)
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	//defer db.Close() // clean up when we're done

	// clientOptions := options.Client().ApplyURI(psn)

	// client, err := mongo.Connect(context.TODO(), clientOptions)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = client.Ping(context.TODO(), nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// db := client.Database(database)
	return _db, nil
}

func GetDbSession(c echo.Context) *mgo.Database {
	_db := c.Get("db").(*mgo.Session)
	return _db.DB(g_database)

}

func GetCollection(c echo.Context, col Collection) (*mgo.Collection, *mgo.Database) {
	_db := c.Get("db").(*mgo.Session)
	colc := string(col)
	return _db.DB(g_database).C(colc), _db.DB(g_database)

}

func GetDBCollection(_db *mgo.Session, col Collection) (*mgo.Collection, *mgo.Database) {
	colc := string(col)
	return _db.DB(g_database).C(colc), _db.DB(g_database)
}
