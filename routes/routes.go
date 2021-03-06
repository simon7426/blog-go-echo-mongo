package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/simon7426/blog-go-echo-mongo/controllers"
)

func BlogRoutes(e *echo.Echo) {
	e.GET("/alive", controllers.Alive)
	// Blogs CRUD
	e.GET("/blogs", controllers.GetAllBlogs)
	e.POST("/blogs", controllers.AddBlog)
	e.GET("/blogs/:blogId", controllers.GetABlog)
	e.PUT("/blogs/:blogId", controllers.EditABlog)
	e.DELETE("/blogs/:blogId", controllers.DeleteABlog)
	// Comments Create & Delete
	e.POST("/blogs/:blogId/comments", controllers.AddAComment)
	e.DELETE("/blogs/:blogId/comments/:commentId", controllers.DeleteAComment)
}
