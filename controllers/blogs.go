package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/simon7426/blog-go-echo-mongo/configs"
	"github.com/simon7426/blog-go-echo-mongo/models"
	"github.com/simon7426/blog-go-echo-mongo/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var blogCollection *mongo.Collection = configs.GetCollection(configs.DB, "blogs")
var validate = validator.New()

func GetAllBlogs(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var blogs []models.Blog

	results, err := blogCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.BlogResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleBlog models.Blog
		if err = results.Decode(&singleBlog); err != nil {
			return c.JSON(http.StatusInternalServerError, response.BlogResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: &echo.Map{
					"data": err.Error(),
				},
			})
		}

		blogs = append(blogs, singleBlog)

	}
	return c.JSON(http.StatusOK, response.BlogResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": blogs,
		},
	})
}

func AddBlog(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var blog models.Blog

	if err := c.Bind(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	if validationErr := validate.Struct(&blog); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": validationErr.Error(),
			},
		})
	}

	newBlog := models.Blog{
		Id:        primitive.NewObjectID(),
		Title:     blog.Title,
		Body:      blog.Body,
		Tags:      blog.Tags,
		Comments:  []models.Comment{},
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}

	result, err := blogCollection.InsertOne(ctx, newBlog)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.BlogResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	return c.JSON(http.StatusCreated, response.BlogResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"data": result,
		},
	})
}

func GetABlog(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("blogId")
	var blog models.Blog

	objID, _ := primitive.ObjectIDFromHex(blogId)

	err := blogCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&blog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.BlogResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, response.BlogResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": blog,
		},
	})
}

func EditABlog(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("blogId")

	var blog models.Blog

	objectId, _ := primitive.ObjectIDFromHex(blogId)

	if err := c.Bind(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	if validationErr := validate.Struct(&blog); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": validationErr.Error(),
			},
		})
	}

	update := bson.M{
		"title":      blog.Title,
		"body":       blog.Body,
		"tags":       blog.Tags,
		"updated_on": time.Now(),
	}

	result, err := blogCollection.UpdateOne(
		ctx,
		bson.M{"id": objectId},
		bson.M{"$set": update},
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.BlogResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	var updatedBlog models.Blog
	if result.MatchedCount == 1 {
		err := blogCollection.FindOne(ctx, bson.M{
			"id": blogId,
		}).Decode(&updatedBlog)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.BlogResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: &echo.Map{
					"data": err.Error(),
				},
			})
		}
	}

	return c.JSON(http.StatusOK, response.BlogResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": updatedBlog,
		},
	})

}

func DeleteABlog(c echo.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("blogId")

	objectId, _ := primitive.ObjectIDFromHex(blogId)

	result, err := blogCollection.DeleteOne(ctx, bson.M{
		"id": objectId,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.BlogResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, response.BlogResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data: &echo.Map{
				"data": "Specified blog not found.",
			},
		})
	}

	return c.JSON(http.StatusOK, response.BlogResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": "User successfully deleted!",
		},
	})
}
