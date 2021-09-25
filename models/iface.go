package models


type RepositoryIface interface {
	FindOneByID(id string) (res Admin, err error)
	InsertOne(admin Admin) error
}
