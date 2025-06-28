package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// NewRequestIDMiddleware creates a new request ID middleware
func NewRequestIDMiddleware() gin.HandlerFunc {
	return requestid.New()
}
