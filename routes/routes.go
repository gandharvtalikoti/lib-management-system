package routes

import (
	"library-management/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r:= gin.Default()

	book:= r.Group("/books")
	{
		book.POST("", controllers.CreateBook)
		book.GET("", controllers.GetBooks)
	}
	return r
}