package database

import (
)


type CollectionInterface interface {
	InsertOne()
	UpdateById()
	Find()
}