package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int
	Message string
	Values  map[string]any
}

type Option func(r *Response)

func New(code int, message string, options ...Option) Response {
	r := Response{
		Code:    code,
		Message: message,
	}

	for _, opt := range options {
		opt(&r)
	}

	return r
}

func WithValues(values map[string]any) Option {
	return func(r *Response) {
		r.Values = values
	}
}

func (r Response) Error(c *gin.Context) {
	c.JSON(r.Code, gin.H{"error": r.Message})
}

func (r Response) OK(c *gin.Context) {
	h := gin.H{}

	if r.Message != "" {
		h["message"] = r.Message
	}

	for k, v := range r.Values {
		h[k] = v
	}

	c.JSON(r.Code, h)
}

func BadRequest(c *gin.Context, err error) {
	New(http.StatusBadRequest, err.Error()).Error(c)
}

func InternalError(c *gin.Context, err error) {
	New(http.StatusInternalServerError, err.Error()).Error(c)
}
