package database

type dbIface interface{
	Collection() collIface
}
type collIface interface{
	FindOne()
}