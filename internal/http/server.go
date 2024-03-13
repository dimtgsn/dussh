package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

const (
	apiPrefix  = "api"
	apiVersion = "v1"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}

func NewServer(
	addr string,
	router http.Handler,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func SetAPIPath(engine *gin.Engine) (*gin.RouterGroup, error) {
	apiPath, err := url.JoinPath(apiPrefix, apiVersion)
	if err != nil {
		return nil, err
	}

	return engine.Group(apiPath), nil
}
