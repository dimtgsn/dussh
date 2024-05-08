package models

import (
	"github.com/gin-gonic/gin"
	"path"
)

const (
	apiPrefix  = "api"
	apiVersion = "v1"
)

var APIPath = makeApiPath()

type Route struct {
	Method   string
	Path     string
	Role     string
	Handlers []gin.HandlerFunc
}

func makeApiPath() string {
	return path.Join(apiPrefix, apiVersion)
}
