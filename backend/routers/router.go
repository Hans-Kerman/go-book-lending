package routers

import (
	"github.com/Hans-Kerman/go-book-lending/backend/controllers"
	"github.com/Hans-Kerman/go-book-lending/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//路由组1:公共访问
	public := r.Group("/api/public")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)

		public.GET("/books", controllers.GetBooksByPage)
		public.GET("/book/:isbn", controllers.GetBookByURL)
	}

	//路由组2:jwt鉴权登录访问
	auth := r.Group("/api")
	{
		auth.Use(middlewares.ParseJWT())

		auth.POST("/borrow", controllers.LendBook)
		auth.POST("/return", controllers.ReturnBook)
		auth.GET("/user/borrows", controllers.GetUserRecord)

		//嵌套路由组3:admin鉴权访问
		admin := auth.Group("/admin")
		{
			admin.Use(middlewares.CheckAdminRole())

			admin.POST("/book", controllers.PostNewBook)
			admin.PUT("/book", controllers.UpdateBook)
		}
	}

	return r
}
