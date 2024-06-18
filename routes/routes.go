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
		bookRoutes.GET("/isbn/:isbn", controllers.GetBookByISBN)
		bookRoutes.GET("/id/:id", controllers.GetBookByID)
		bookRoutes.POST("", controllers.CreateBook)
		bookRoutes.GET("", controllers.GetBooks)
	}

	issuedBookRoutes:= r.Group("/issued")
	{
		issuedBookRoutes.POST("", controllers.IssueBook)
		issuedBookRoutes.GET("/:user_id", controllers.GetIssuedBooksByUser)
		issuedBookRoutes.GET("overdue", controllers.GetOverdueBooks)
		issuedBookRoutes.GET("overdue/:user_id", controllers.GetOverdueBooksByUser)
	}

	r.GET("/search", controllers.SearchBook)
	return r
}
