package main

import (
	"flag"
	"fmt"

	"github.com/masagatech/myinvest/module"

	"github.com/masagatech/myinvest/db"
	"gopkg.in/mgo.v2"

	"github.com/go-playground/validator"

	"github.com/masagatech/myinvest/api"
	"github.com/masagatech/myinvest/server"
	"github.com/masagatech/myinvest/util"
)

func main() {

	cfgPath := flag.String("p", "./conf.local.yaml", "Path to config file")
	flag.Parse()
	cfg, err := util.Load(*cfgPath)

	checkErr(err)

	db, err := sartdb(cfg)

	// onload brokers
	module.Brokers_init(db)
	checkErr(err)
	checkErr(start(cfg, db))

}

// func middleware(cfg *util.Configuration) error {

// }

func sartdb(cfg *util.Configuration) (*mgo.Session, error) {
	fmt.Println(cfg.DB.URL)
	db, err := db.New(cfg.DB.URL, cfg.DB.DbName)
	if err != nil {
		return nil, err
	}
	return db, err

}

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func start(cfg *util.Configuration, db *mgo.Session) error {
	e := server.New(cfg, db)

	//JWT TOKEN
	e.Validator = &CustomValidator{validator: validator.New()}
	api.RegisterRoute(cfg, e)
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil

}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
