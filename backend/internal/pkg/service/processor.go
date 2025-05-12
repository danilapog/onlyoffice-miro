package service

type StorageProcessor[ID comparable, T any, R any] interface {
	TableName() string
	BuildSelectQuery(id ID) (query string, args []any, scanner func(R) (T, error))
	BuildInsertQuery(id ID, component T) (query string, args []any)
	BuildUpdateQuery(id ID, component T) (query string, args []any)
	BuildDeleteQuery(id ID) (query string, args []any)
}
