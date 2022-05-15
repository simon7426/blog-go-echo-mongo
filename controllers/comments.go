package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/simon7426/blog-go-echo-mongo/models"
	"github.com/simon7426/blog-go-echo-mongo/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAComment(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var comment models.Comment

	blogId := c.Param("blogId")
	objectId, _ := primitive.ObjectIDFromHex(blogId)

	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	if validationErr := validate.Struct(&comment); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": validationErr.Error(),
			},
		})
	}

	if comment.DisplayName == "" {
		comment.DisplayName = comment.Email
	}

	newComment := models.Comment{
		Id:          primitive.NewObjectID(),
		Email:       comment.Email,
		DisplayName: comment.DisplayName,
		Body:        comment.Body,
		CreatedOn:   time.Now(),
	}

	result, err := blogCollection.UpdateOne(
		ctx,
		bson.M{"id": objectId},
		bson.M{"$push": bson.M{
			"comments": newComment,
		}},
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	var updatedBlog models.Blog
	if result.ModifiedCount == 1 {
		err := blogCollection.FindOne(ctx, bson.M{
			"id": objectId,
		}).Decode(&updatedBlog)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.BlogResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: &echo.Map{
					"data": err.Error(),
				},
			})
		}
	}

	if result.ModifiedCount == 0 {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": fmt.Sprintf("Unable to add Comment in %s", blogId),
			},
		})
	}

	return c.JSON(http.StatusOK, response.BlogResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": updatedBlog,
		},
	})

}

func DeleteAComment(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("blogId")
	commentId := c.Param("commentId")

	blogObjectId, _ := primitive.ObjectIDFromHex(blogId)
	commentObjectId, _ := primitive.ObjectIDFromHex(commentId)

	result, err := blogCollection.UpdateOne(
		ctx,
		bson.M{"id": blogObjectId},
		bson.M{"$pull": bson.M{
			"comments": bson.M{"id": commentObjectId},
		},
		},
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": err.Error(),
			},
		})
	}

	var updatedBlog models.Blog
	if result.ModifiedCount == 1 {
		err := blogCollection.FindOne(ctx, bson.M{
			"id": blogObjectId,
		}).Decode((&updatedBlog))
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.BlogResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: &echo.Map{
					"data": err.Error(),
				},
			})
		}
	}

	if result.ModifiedCount == 0 {
		return c.JSON(http.StatusBadRequest, response.BlogResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &echo.Map{
				"data": fmt.Sprintf("Comment %s not found", commentId),
			},
		})
	}

	return c.JSON(http.StatusOK, response.BlogResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": updatedBlog,
		},
	})

}
