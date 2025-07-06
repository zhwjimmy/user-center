package middleware

import "github.com/gin-gonic/gin"

type (
	RecoveryMiddleware  gin.HandlerFunc
	LoggerMiddleware    gin.HandlerFunc
	RequestIDMiddleware gin.HandlerFunc
	CORSMiddleware      gin.HandlerFunc
)
