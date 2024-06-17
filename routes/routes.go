package routes

import (
	"library-management/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r:= gin.Default()

	userRoutes:= r.Group("/users")
	{
		userRoutes.POST("", controllers.CreateUser)
		userRoutes.GET("", controllers.GetUsers)
	}
	bookRoutes:= r.Group("/books")
	{
		bookRoutes.GET("/:isbn", controllers.GetBookByISBN)
		bookRoutes.POST("", controllers.CreateBook)
		bookRoutes.GET("", controllers.GetBooks)
	}

	issuedBookRoutes:= r.Group("/issued")
	{
		issuedBookRoutes.POST("", controllers.IssueBook)
	}
	return r
}
