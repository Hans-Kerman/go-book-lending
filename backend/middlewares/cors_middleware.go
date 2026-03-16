package middlewares

import (
	"github.com/gin-contrib/cors"
)

var CORSConf = cors.Config{
	AllowOrigins:     []string{"http://localhost:5173"},
	AllowMethods:     []string{"GET", "POST"},
	AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * 3600,
}
