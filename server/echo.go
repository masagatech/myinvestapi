package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/masagatech/myinvest/util"
	"gopkg.in/mgo.v2"
)

func New(cfg *util.Configuration, db *mgo.Session) *echo.Echo {
	e := echo.New()

	//e.Use(session.Middleware(sessions.NewCookieStore([]byte("sec2"))))
	e.Use(redisMiddleware(redisInstance(cfg)))
	e.Use(session.Middleware(sessions.NewFilesystemStore("", []byte("key-1"), nil)))
	e.Use(middleware.Logger(), middleware.Recover(),
		CORS(), Headers(),
		dataSourceMiddleware(db, cfg),
	)

	e.GET("/", healthCheck)
	//e.Validator = &CustomValidator{V: validator.New()}
	//custErr := &customErrHandler{e: e}
	//e.HTTPErrorHandler = custErr.handler
	//e.Binder = &CustomBinder{b: &echo.DefaultBinder{}}
	return e
}

func redisInstance(cfg *util.Configuration) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password, // no password set
		DB:       cfg.Redis.DB,       // use default DB
	})

	return rdb

}

func redisMiddleware(_redis *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			ectx.Set("redis", _redis)
			return next(ectx)
		}
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func dataSourceMiddleware(mongo *mgo.Session, cnf *util.Configuration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			dbsession := mongo.Copy()
			defer dbsession.Close() // clean up
			fmt.Println("inside adapter")
			c.Set("db", dbsession)
			c.Set("config", cnf)

			return next(c)
		}
	}
}

// Headers adds general security headers for basic security measures
func Headers() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Protects from MimeType Sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			// Prevents browser from prefetching DNS
			c.Response().Header().Set("X-DNS-Prefetch-Control", "off")
			// Denies website content to be served in an iframe
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
			// Prevents Internet Explorer from executing downlo	ads in site's context
			c.Response().Header().Set("X-Download-Options", "noopen")
			// Minimal XSS protection
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			return next(c)
		}
	}
}

// CORS adds Cross-Origin Resource Sharing support
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		MaxAge:           86400,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

// Config represents server specific config
type Config struct {
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	Debug               bool
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// Start starts echo server
func Start(e *echo.Echo, cfg *Config) {
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	e.Debug = cfg.Debug

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./template/*.html")),
	}
	e.Renderer = renderer

	// Start server
	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
