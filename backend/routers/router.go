package routers

import (
	"github.com/Hans-Kerman/go-book-lending/backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//路由组1:公共访问
	{
		public := r.Group("/api")

		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)

		public.GET("/books", controllers.GetBooksByPage)
		public.GET("/book/:isbn", controllers.GetBookByISBN)
	}

	return r
}
