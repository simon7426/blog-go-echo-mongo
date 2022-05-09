package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	Tags      []string           `json:"tags"`
	Comments  []Comment          `json:"comments,omitempty"`
	CreatedOn time.Time          `json:"created_on,omitempty"`
	UpdatedOn time.Time          `json:"updated_on,omitempty"`
}

type BlogCollection struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Title     string             `json:"title"`
	Tags      []string           `json:"tags,omitempty"`
	Link      string             `json:"link,omitempty"`
	CreatedOn time.Time          `json:"created_on,omitempty"`
	UpdatedOn time.Time          `json:"updated_on,omitempty"`
}

type Comment struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Email       string             `json:"email"`
	DisplayName string             `json:"display_name,omitempty"`
	Body        string             `json:"body"`
	CreatedOn   time.Time          `json:"created_on,omitempty"`
}
