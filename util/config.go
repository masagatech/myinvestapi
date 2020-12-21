package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/masagatech/myinvest/model"
	"gopkg.in/yaml.v2"
)

var (
	Config *Configuration
)

// Load returns Configuration struct
func Load(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}
	var cfg = new(Configuration)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	Config = cfg

	return cfg, nil
}

// Configuration holds data necessary for configuring application
type Configuration struct {
	Server  *Server      `yaml:"server,omitempty"`
	Zerodha *Zerodha     `yamal:"zerodha,omitempty"`
	DB      *Database    `yaml:"mongo,omitempty"`
	JWT     *JWT         `yaml:"jwt,omitempty"`
	App     *Application `yaml:"application,omitempty"`
}

// Database holds data necessary for database configuration
type Database struct {
	LogQueries bool   `yaml:"log_queries,omitempty"`
	Timeout    int    `yaml:"timeout_seconds,omitempty"`
	URL        string `yaml:"url,omitempty"`
	DbName     string `yaml:"db,omitempty"`
}

// Server holds data necessary for server configuration
type Server struct {
	Port         string `yaml:"port,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
}

// JWT holds data necessary for JWT configuration
type JWT struct {
	MinSecretLength  int    `yaml:"min_secret_length,omitempty"`
	DurationMinutes  int    `yaml:"duration_minutes,omitempty"`
	RefreshDuration  int    `yaml:"refresh_duration_minutes,omitempty"`
	MaxRefresh       int    `yaml:"max_refresh_minutes,omitempty"`
	SigningAlgorithm string `yaml:"signing_algorithm,omitempty"`
	Secret           string `yaml:"secret,omitempty"`
}

// Application holds application configuration details
type Application struct {
	MinPasswordStr int    `yaml:"min_password_strength,omitempty"`
	SwaggerUIPath  string `yaml:"swagger_ui_path,omitempty"`
}

type ResponseModel struct {
	ErrorCode   string      `json:"errorCode"`
	Message     string      `json:"message"`
	ResultKey   int         `json:"resultKey"`
	ResultValue interface{} `json:"resultValue"`
}

type Zerodha struct {
	ApiUrl    string `yaml:"api_url,omitempty"`
	ApiKey    string `yaml:"api_key,omitempty"`
	ApiSecret string `yaml:"api_secret,omitempty"`
}

func BindSysParams(cmp string, operation string, flag string, payload string) (string, model.Sysparams) {
	var sysScema model.Sysparams
	sysScema.Schema = "cmp" + cmp
	sysScema.Flag = flag
	sysScema.Operate = operation
	sysScema.Payload = payload
	byteArray, err := json.Marshal(sysScema)

	if err != nil {
		return "{}", sysScema
	}

	return string(byteArray), sysScema
}

func ToJsonString(params interface{}) string {

	byteArray, err := json.Marshal(params)

	if err != nil {
		return "{}"
	}

	return string(byteArray)
}

func CreateResponse(c echo.Context, resultKey int, data interface{}, errorCode string) error {

	res := ResponseModel{
		ResultKey:   resultKey,
		ErrorCode:   errorCode,
		ResultValue: data,
	}
	if resultKey == 1 {
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, res)

}

func GetConfig(c echo.Context) *Configuration {
	return c.Get("config").(*Configuration)
}
